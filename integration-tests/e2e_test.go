package integrationtests

import (
	"encoding/gob"
	"net"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/helpers"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/services/sstorage"
	"github.com/tudorhulban/analytics77/shared"
)

func TestAnalytics_E2E(t *testing.T) {
	localTime := time.Now()
	_, offsetSeconds := localTime.Zone()

	utcOffsetHours := int64(offsetSeconds / 3600)

	utcOffset := helpers.TimestampOffsets{
		OffsetUTCHours: utcOffsetHours,
	}

	apiKey := os.Getenv(cmd.OSAPIGeolocation)

	serviceLogging, fnCloseLogging, erCrServiceLogging := slogging.NewServiceLogging("", os.Stdout)
	require.NoError(t, erCrServiceLogging)

	defer fnCloseLogging()

	serviceStorage, errCrServiceStorage := sstorage.NewServiceStorage(t.TempDir())
	require.NoError(t, errCrServiceStorage)
	require.NotNil(t, serviceStorage)

	t.Cleanup(
		func() {
			require.NoError(t, serviceStorage.Close())
		},
	)

	serviceGeo, errCrServiceGeo := sgeo.NewServiceGeo(
		&sgeo.ParamsNewServiceGeo{
			APIKeyGeolocation: apiKey,
		},
		serviceStorage,
	)
	require.NoError(t, errCrServiceGeo)
	require.NotNil(t, serviceGeo)

	dummyURL, _ := url.Parse("https://example.com/analytics")

	tests := []struct {
		description   string
		inputRequests shared.Requests

		expectedRecordCount int
		expectedSitesCount  int
	}{
		{
			description: "1. Send single request",
			inputRequests: shared.Requests{
				{
					RemoteAddr: "82.77.237.37",
					Host:       "example.com",
					Method:     "POST",
					URL:        dummyURL,
					Header:     map[string][]string{"Content-Type": {"application/json"}},

					TimestampUNIX:  localTime.Unix(),
					OffsetUTCHours: utcOffsetHours,
				},
			},
			expectedRecordCount: 1,
			expectedSitesCount:  1,
		},
		{
			description: "2. Send multiple requests in one batch",
			inputRequests: shared.Requests{
				{
					RemoteAddr: "82.77.237.38",
					Host:       "api.com",
					Method:     "GET",
					URL:        dummyURL,

					TimestampUNIX:  localTime.Unix(),
					OffsetUTCHours: utcOffsetHours,
				},
				{
					RemoteAddr: "82.77.237.39",
					Host:       "metrics.com",
					Method:     "PUT",
					URL:        dummyURL,

					TimestampUNIX:  localTime.Unix(),
					OffsetUTCHours: utcOffsetHours,
				},
			},
			expectedRecordCount: 2,
			expectedSitesCount:  2,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				listener, errListener := net.Listen("tcp", "127.0.0.1:0") //nolint:noctx
				require.NoError(t, errListener)

				serviceAnalytics, errCrServiceAnalytics := sanalytics.NewServiceAnalytics(
					&sanalytics.PiersNewServiceAnalytics{
						ServiceGeo: serviceGeo,
					},
					&utcOffset,
				)
				require.NoError(t, errCrServiceAnalytics)
				require.NotNil(t, serviceAnalytics)

				transportTCP, errCrTransport := transporttcp.NewTransportTCP(
					listener,
					&transporttcp.PiersNewTransportTCP{
						ServiceLogging:   serviceLogging,
						ServiceAnalytics: serviceAnalytics,
					},
				)
				require.NoError(t, errCrTransport)

				go func() {
					if errServerStart := transportTCP.Start(); errServerStart != nil {
						serviceLogging.Logger.Printf(
							"transport TCP stopped: %v",
							errServerStart,
						)
					}
				}()

				// Give the OS a tiny moment to bind the socket
				time.Sleep(10 * time.Millisecond)

				connClient, errListener := net.Dial( //nolint:noctx
					"tcp",
					transportTCP.GetListeningAddress(),
				)
				require.NoError(t, errListener)
				require.NotZero(t, connClient)

				require.NoError(t,
					gob.NewEncoder(connClient).Encode(&tc.inputRequests),
				)

				// Close to flush and trigger EOF on the server side.
				connClient.Close()

				// wrapping up and waiting for processing to stop.
				transportTCP.Stop()

				time.Sleep(1 * time.Second)

				// compare site names
				require.EqualValues(t,
					tc.expectedSitesCount,
					len(serviceAnalytics.DC.GetSiteNames()),

					"comparing number of site names",
				)

				hoursWithData, noHours := serviceAnalytics.DC.CurrentDayHoursWithData(&utcOffset)
				require.NotZero(t, noHours)
				require.LessOrEqual(t, noHours, int8(2))
				require.NotZero(t, hoursWithData)
			},
		)
	}
}
