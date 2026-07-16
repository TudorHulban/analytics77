package shared

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/domain/analytics"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/sstorage"
)

func ForwardAnalytics(r *http.Request, encoder *gob.Encoder) error {
	flat := Request{
		RemoteAddr: r.RemoteAddr,
		Host:       r.Host,
		Method:     r.Method,
		URL:        r.URL,
		Header:     r.Header,
	}

	// This streams the raw binary directly across the TCP socket
	return encoder.Encode(flat)
}

func TestRequest_AsParamsAddEvent(t *testing.T) {
	serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(t.TempDir())
	require.NoError(t, errCrServiceStorage)
	require.NotNil(t, serviceStorage)

	t.Cleanup(
		func() {
			require.NoError(t, serviceStorage.Close())
		},
	)

	apiKey := os.Getenv(cmd.OSAPIGeolocation)

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
		&sgeo.ParamsNewServiceGeo{
			APIKeyGeolocation: apiKey,
		},
		serviceStorage,
	)
	require.NoError(t, errCrServiceGeo)
	require.NotNil(t, serviceGeo)

	offsetsROU := helpers.TimestampOffsets{
		OffsetUTCHours: -3,
	}

	tests := []struct {
		piers    *PiersAsParamsAddEvent
		validate func(t *testing.T, res *ParamsAddEvent)

		name        string
		errContains string
		req         Request

		wantErr bool
	}{
		{
			name: "1. Error - Offsets is nil",
			req:  Request{RemoteAddr: "127.0.0.1:8080"},
			piers: &PiersAsParamsAddEvent{
				Offsets:    nil,
				ServiceGeo: serviceGeo,
			},
			wantErr:     true,
			errContains: "AsParamsAddEvent - passed offsets is nil",
		},
		{
			name: "2. Error - ServiceGeo is nil",
			req:  Request{RemoteAddr: "127.0.0.1:8080"},
			piers: &PiersAsParamsAddEvent{
				Offsets:    &offsetsROU,
				ServiceGeo: nil,
			},
			wantErr:     true,
			errContains: "AsParamsAddEvent - passed ServiceGeo is nil",
		},
		{
			name: "3. Error - Invalid IP Address",
			req:  Request{RemoteAddr: "not-an-ip-address"},
			piers: &PiersAsParamsAddEvent{
				Offsets:    &offsetsROU,
				ServiceGeo: serviceGeo,
			},
			wantErr: true, // Should fail at netip.ParseAddr
		},
		{
			name: "4. Happy Path - Chrome User Agent and Valid Parse",
			req: Request{
				RemoteAddr:     "82.77.237.37:53",
				Host:           "some-site.eu",
				Header:         map[string][]string{"User-Agent": {"Mozilla/5.0 Chrome/120.0.0.0"}},
				TimestampUNIX:  1717500000,
				OffsetUTCHours: 2,
			},
			piers: &PiersAsParamsAddEvent{
				Offsets:    &offsetsROU,
				ServiceGeo: serviceGeo,
			},
			wantErr: false,
			validate: func(t *testing.T, res *ParamsAddEvent) {
				t.Helper()

				assert.Equal(t,
					"some-site.eu",
					res.SiteKey,
				)
				assert.NotZero(t, res.Country)
				assert.NotZero(t, res.City)

				assert.Equal(t,
					"82.77.237.37",
					res.IP,
				)
				assert.Equal(t, analytics.Chrome, res.Browser)
				assert.NotZero(t, res.ASNOrganization)

				require.EqualValues(t,
					1717500000,
					res.TimestampUNIX,
					"timestamp mismtch",
				)
				require.EqualValues(t,
					2,
					res.OffsetUTCHours,
					"offset mismatch",
				)
			},
		},
		{
			name: "5. Happy Path - Safari User Agent Fallback (No Chrome)",
			req: Request{
				RemoteAddr: "82.77.237.38",
				Header:     map[string][]string{"User-Agent": {"Mozilla/5.0 Safari/605.1.15"}},
			},
			piers: &PiersAsParamsAddEvent{
				Offsets:    &offsetsROU,
				ServiceGeo: serviceGeo,
			},
			wantErr: false,
			validate: func(t *testing.T, res *ParamsAddEvent) {
				t.Helper()

				assert.Equal(t,
					"82.77.237.38",
					res.SiteKey,
					"should failback to IP",
				)
				assert.Equal(t, analytics.Safari, res.Browser)
			},
		},
		{
			name: "6. Happy Path - Unknown Browser Default",
			req: Request{
				RemoteAddr: "127.0.0.1",
				Header:     map[string][]string{"User-Agent": {"Curl/7.81.0"}},
			},
			piers: &PiersAsParamsAddEvent{
				Offsets:    &offsetsROU,
				ServiceGeo: serviceGeo,
			},
			wantErr: false,
			validate: func(t *testing.T, res *ParamsAddEvent) {
				t.Helper()

				assert.Equal(t, analytics.Browser(0), res.Browser)
			},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				result, errTransformation := tc.req.AsParamsAddEvent(tc.piers)

				if tc.wantErr {
					require.Error(t, errTransformation)

					if tc.errContains != "" {
						assert.Contains(t, errTransformation.Error(), tc.errContains)
					}

					assert.Nil(t, result)

					return
				}

				require.NoError(t, errTransformation)
				require.NotNil(t, result)

				if tc.validate != nil {
					tc.validate(t, result)
				}
			},
		)
	}
}
