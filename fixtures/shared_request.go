package fixtures

import (
	"net/http"
	"net/url"
	"time"

	"github.com/tudorhulban/analytics77/shared"
)

func NewRequest(withHost, withIP string) shared.Request {
	now := time.Now()
	_, offsetSecs := now.Zone()

	return shared.Request{
		// Fallback address (with dummy port so net.SplitHostPort does not fail)
		RemoteAddr: withIP + ":12345",
		Host:       withHost,
		Method:     http.MethodGet,

		// We just instantiate an empty struct pointer to satisfy *url.URL
		// without needing url.Parse() or error handling.
		URL: &url.URL{
			Host: withHost,
		},

		Header: map[string][]string{
			"X-Forwarded-For": {withIP},
			"X-Real-IP":       {withIP},
		},

		TimestampUNIX:  now.Unix(),
		OffsetUTCHours: int64(offsetSecs / 3600),
	}
}

func NewRequests(withHost string, withIPs ...string) shared.Requests {
	result := make([]shared.Request, len(withIPs))

	for ix, withIP := range withIPs {
		result[ix] = NewRequest(withHost, withIP)
	}

	return result
}
