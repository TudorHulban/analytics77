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
	"github.com/tudorhulban/analytics77/domain/offsets"
	"github.com/tudorhulban/analytics77/domain/site"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/analytics77/services/sstorage"
	"github.com/tudorhulban/analytics77/shared"
)

func TestTransport_TCP(t *testing.T) {
	t.Parallel()

	dummyURL, _ := url.Parse("https://example.com/analytics")

	offsetsROU := offsets.TimestampOffsets{
		OffsetUTCHours: -3,
	}

	serviceLogging, fnCloseLogging, errCrServiceLogging := slogging.NewServiceLogging("", os.Stdout)
	require.NoError(t, errCrServiceLogging)

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
			APIKeyGeolocation: os.Getenv(cmd.OSAPIGeolocation),
		},
		serviceStorage,
	)
	require.NoError(t, errCrServiceGeo)
	require.NotNil(t, serviceGeo)

	tests := []struct {
		description   string
		inputRequests shared.Requests

		expectedSites          []site.SiteName
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

			expectedCountSites:     1,
			expectedSites:          []site.SiteName{"example.com"},
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

			expectedCountSites:     2,
			expectedSites:          []site.SiteName{"api.com", "metrics.com"},
			expectedRecordsPerSite: []uint32{1, 1},
		},
	}

	for _, testCase := range tests {
		t.Run(
			testCase.description,
			func(t *testing.T) {
				t.Parallel()

				listener, errListener := net.Listen("tcp", "127.0.0.1:0") //nolint:noctx
				require.NoError(t, errListener)

				dc, errCr := datacenter.NewDataCenter()
				require.NoError(t, errCr)

				transportTCP, errCrTransport := NewTransportTCP(
					listener,
					&PiersNewTransportTCP{
						ServiceLogging: serviceLogging,
						ServiceGeo:     serviceGeo,
						Datacenter:     dc,
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
					gob.NewEncoder(connClient).Encode(&testCase.inputRequests),
				)

				// Close to flush and trigger EOF on the server side.
				connClient.Close()

				// wrapping up and waiting for processing to stop.
				transportTCP.Stop()

				require.EqualValues(t,
					testCase.expectedCountSites,
					len(transportTCP.Datacenter.GetSiteNames()),
				)

				hoursWithData, howMany := transportTCP.
					Datacenter.
					CurrentDayHoursWithData(&offsetsROU)
				require.EqualValues(t, 1, howMany)

				fmt.Println(hoursWithData)

				recordsPerSite := transportTCP.
					Datacenter.
					GetCurrentHourRecordsPerSite(&offsetsROU)

				transportTCP.Datacenter.Snapshot(os.Stdout)

				fmt.Println(recordsPerSite)

				recordsPerSite.Verify(t,
					testCase.expectedSites,
					testCase.expectedRecordsPerSite,
				)

				require.EqualValues(t,
					testCase.expectedCountSites,
					len(transportTCP.Datacenter.GetCurrentHourRecordsPerSite(&offsetsROU)),
				)
			},
		)
	}
}
