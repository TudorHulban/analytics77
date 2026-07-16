package sgeo

import (
	"errors"
	"net/http"
	"net/netip"
	"time"

	"github.com/tudorhulban/analytics77/domain/analytics"
	requestgeo "github.com/tudorhulban/analytics77/infra/request-geo"
	"github.com/tudorhulban/analytics77/services/sstorage"
	lru "github.com/tudorhulban/hx-lru"
	"github.com/tudorhulban/hxerrors"
)

// 1. Check LRU
// 2. If hit → return
// 3. If miss → call ServiceStorage.GetIPGeo(ip)
// 4. If storage hit → store in LRU → return
// 5. If storage miss → call geo provider
// 6. Persist provider result into Bitcask
// 7. Store in LRU
// 8. Return

type ServiceGeo struct {
	cache          *lru.CacheOneLRU[string, analytics.GeoIP]
	serviceStorage *sstorage.ServiceStorage
	httpClient     *http.Client

	apiKeyGeolocation string
}

type ParamsNewServiceGeo struct {
	APIKeyGeolocation string
}

func NewServiceGeo(params *ParamsNewServiceGeo, serviceStorage *sstorage.ServiceStorage) (*ServiceGeo, error) {
	if serviceStorage == nil {
		return nil,
			errors.New("passed service storage is nil")
	}

	return &ServiceGeo{
			apiKeyGeolocation: params.APIKeyGeolocation,

			cache: lru.NewCacheOneLRU[string, analytics.GeoIP](
				&lru.ParamsNewCacheLRU{
					TTL:      14 * 24 * time.Hour,
					Capacity: 5000,
				},
			),
			serviceStorage: serviceStorage,
			httpClient: &http.Client{
				Timeout: 5 * time.Second,
			},
		},
		nil
}

func (s *ServiceGeo) GetIPGeo(ip netip.Addr) (*analytics.GeoIP, error) {
	if !ip.IsValid() {
		return nil, hxerrors.ErrNilInput{InputName: "ip"}
	}

	ipStr := ip.String()

	if cacheValue, errGetCache := s.cache.Get(ipStr); errGetCache == nil {
		return cacheValue,
			nil
	}

	if kvValue, errStorage := s.serviceStorage.GetIPGeo(ipStr); errStorage == nil {
		s.cache.Put(ipStr, *kvValue)

		return kvValue, nil
	}

	providerValue, errGetGeolocation := requestgeo.GetLocationByIP(
		&requestgeo.ParamsGetLocationByIP{
			Client:    s.httpClient,
			APIKey:    s.apiKeyGeolocation,
			IPAddress: ipStr,
		},
	)
	if errGetGeolocation != nil {
		return nil,
			errGetGeolocation
	}

	_ = s.serviceStorage.PutGeoIP(providerValue)
	s.cache.Put(ipStr, *providerValue)

	return providerValue, nil
}
