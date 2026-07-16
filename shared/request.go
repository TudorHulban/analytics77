package shared

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"strings"

	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/services/sgeo"

	"github.com/tudorhulban/hxhelpers"
)

type Request struct {
	Header map[string][]string
	URL    *url.URL

	RemoteAddr string
	Host       string
	Method     string

	TimestampUNIX  int64
	OffsetUTCHours int64
}

type PiersAsParamsAddEvent struct {
	Offsets    *helpers.TimestampOffsets
	ServiceGeo *sgeo.ServiceGeo
}

func (req Request) AsParamsAddEvent(dependencies *PiersAsParamsAddEvent) (*ParamsAddEvent, error) {
	if dependencies.Offsets == nil {
		return nil,
			errors.New(
				"AsParamsAddEvent - passed offsets is nil",
			)
	}

	if dependencies.ServiceGeo == nil {
		return nil,
			errors.New(
				"AsParamsAddEvent - passed ServiceGeo is nil",
			)
	}

	host, _, errHost := net.SplitHostPort(req.RemoteAddr)
	if errHost != nil {
		host = req.RemoteAddr
	}

	ip, errParseIP := netip.ParseAddr(host)
	if errParseIP != nil {
		return nil,
			fmt.Errorf(
				"parsing IP: %s: %w",
				ip,
				errParseIP,
			)
	}

	geoInfo := analytics.GeoIP{
		IsPrivate: true,
	}

	if !ip.IsPrivate() && !ip.IsLoopback() {
		responseGeo, errGeo := dependencies.ServiceGeo.GetIPGeo(ip)
		if errGeo != nil {
			return nil,
				fmt.Errorf(
					"geo call for IP: %s: %w",
					ip,
					errGeo,
				)
		}

		geoInfo = *responseGeo
	}

	userAgent := req.Header["User-Agent"]

	var uaString string

	if len(userAgent) > 0 {
		uaString = userAgent[0]
	}

	var browser analytics.Browser

	switch {
	case strings.Contains(uaString, "CriOS"):
		browser = analytics.Chrome

	case strings.Contains(uaString, "Safari") && !strings.Contains(uaString, "Chrome"):
		browser = analytics.Safari

	case strings.Contains(uaString, "Edg"):
		browser = analytics.Edge

	case strings.Contains(uaString, "Firefox"):
		browser = analytics.Firefox

	case strings.Contains(uaString, "Brave"):
		browser = analytics.Brave

	case strings.Contains(uaString, "Chrome"):
		browser = analytics.Chrome

	default:
		browser = 0
	}

	return &ParamsAddEvent{
			SiteKey: hxhelpers.Ternary(
				len(req.Host) == 0,

				host,
				req.Host,
			),
			Country: geoInfo.Location.CountryCode,
			City:    geoInfo.Location.City,

			IP:      ip.String(),
			Browser: browser,

			ASNOrganization: geoInfo.ASN.Organization,

			OffsetUTCHours: hxhelpers.Ternary(
				req.OffsetUTCHours > 0,

				req.OffsetUTCHours,
				dependencies.Offsets.OffsetUTCHours,
			),
			TimestampUNIX: req.TimestampUNIX,

			IsPrivateIP: geoInfo.IsPrivate,
		},
		nil
}

type Requests []Request
