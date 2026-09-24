package sstorage

import (
	"github.com/prologic/bitcask"
	"github.com/shamaton/msgpack/v3"
	"github.com/tudorhulban/analytics77/domain/analytics"
)

type ServiceStorage struct {
	db *bitcask.Bitcask
}

func NewServiceStorage(path string) (*ServiceStorage, error) {
	db, errCrBitcaskDB := bitcask.Open(path)
	if errCrBitcaskDB != nil {
		return nil,
			errCrBitcaskDB
	}

	return &ServiceStorage{
			db: db,
		},
		nil
}

func (s *ServiceStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}

	return nil
}

func (*ServiceStorage) PutGeoIP(value *analytics.GeoIP) error {
	return nil
}

func (s *ServiceStorage) GetIPGeo(ip string) (*analytics.GeoIP, error) {
	key := []byte("ip:" + ip)

	dbValue, errGet := s.db.Get(key)
	if errGet != nil {
		if errGet == bitcask.ErrKeyNotFound {
			return nil,
				ErrIPNotFound
		}

		return nil, errGet
	}

	var result analytics.GeoIP

	if errUnmarshal := msgpack.Unmarshal(dbValue, &result); errUnmarshal != nil {
		return nil,
			errUnmarshal
	}

	return &result, nil
}
