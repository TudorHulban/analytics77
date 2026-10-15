package transporttcp

import (
	"encoding/gob"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/services/sstorage"
	"github.com/tudorhulban/analytics77/shared"
)

func TestTransport_TCP(t *testing.T) {
	t.Parallel()

	dummyURL, _ := url.Parse("https://example.com/analytics")

	offsetsROU := helpers.TimestampOffsets{
		OffsetUTCHours: -3,
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

	tests := []struct {
		description   string
		inputRequests shared.Requests

		expectedSites          []datacenter.Site
		expectedRecordsPerSite []uint32
		expectedCountSites     int
	}{
		{
			description: "1. Send single request",
			inputRequests: shared.Requests{
				{
					RemoteAddr:    "10.0.0.1",
					Host:          "example.com",
					Method:        http.MethodPost,
					URL:           dummyURL,
					Header:        map[string][]string{"Content-Type": {"application/json"}},
					TimestampUNIX: time.Now().Unix(),
				},
			},

			expectedCountSites: 1,

			expectedSites:          []datacenter.Site{"example.com"},
			expectedRecordsPerSite: []uint32{1},
		},
		{
			description: "2. Send multiple requests in one batch",
			inputRequests: shared.Requests{
				{
					RemoteAddr:    "10.0.0.2",
					Host:          "api.com",
					Method:        "GET",
					URL:           dummyURL,
					TimestampUNIX: time.Now().Unix(),
				},
				{
					RemoteAddr:    "10.0.0.3",
					Host:          "metrics.com",
					Method:        "PUT",
					URL:           dummyURL,
					TimestampUNIX: time.Now().Unix(),
				},
			},
			expectedCountSites: 2,

			expectedSites:          []datacenter.Site{"api.com", "metrics.com"},
			expectedRecordsPerSite: []uint32{1, 1},
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.description,
			func(t *testing.T) {
				t.Parallel()

				listener, errListener := net.Listen("tcp", "127.0.0.1:0") //nolint:noctx
				require.NoError(t, errListener)

				serviceAnalytics, errCrAnalytics := sanalytics.NewServiceAnalytics(
					&sanalytics.PiersNewServiceAnalytics{
						ServiceGeo: serviceGeo,
					},
					&offsetsROU,
				)
				require.NoError(t, errCrAnalytics)

				transportTCP, errCrTransport := NewTransportTCP(
					listener,
					&PiersNewTransportTCP{
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
					transportTCP.listener.Addr().String(),
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

				require.EqualValues(t,
					tc.expectedCountSites,
					len(serviceAnalytics.DC.GetSiteNames()),
				)

				hoursWithData, howMany := serviceAnalytics.DC.CurrentDayHoursWithData(&offsetsROU)
				require.EqualValues(t, 1, howMany)

				fmt.Println(hoursWithData)

				recordsPerSite := serviceAnalytics.DC.GetCurrentHourRecordsPerSite(&offsetsROU)

				serviceAnalytics.DC.Snapshot(os.Stdout)

				fmt.Println(recordsPerSite)

				require.True(t,
					recordsPerSite.Verify(
						tc.expectedSites,
						tc.expectedRecordsPerSite,
					),
				)

				require.EqualValues(t,
					tc.expectedCountSites,
					len(serviceAnalytics.DC.GetCurrentHourRecordsPerSite(&offsetsROU)),
				)
			},
		)
	}
}
