package appanalytics

import (
	"fmt"
	"net"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/helpers"
	"github.com/tudorhulban/analytics77/infra/initialization"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/slogging"
	"github.com/tudorhulban/hxerrors"
)

type ParamsInitializeApp struct {
	ConfigPortRPC  string
	ConfigPortHTTP string

	KeyGeolocationAPI string
	PathLogFile       string
}

func InitializeApp(params *ParamsInitializeApp) *App {
	listener, errListener := net.Listen( //nolint:noctx
		"tcp",
		fmt.Sprintf(
			"127.0.0.1:%s",
			params.ConfigPortRPC,
		),
	)
	if errListener != nil {
		fmt.Printf(
			"error create listener: %s\n",
			errListener.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConnectivityIssues,
		)
	}

	serviceLogging, fnCloseLogging, erCrServiceLogging := slogging.NewServiceLogging(params.PathLogFile, os.Stdout)
	if erCrServiceLogging != nil {
		fmt.Printf(
			"error create servce logging: %s\n",
			erCrServiceLogging.Error(),
		)

		os.Exit(
			hxerrors.OSExitForLoggingIssues,
		)
	}

	serviceAnalytics, errInitialization := initialization.Services(
		&initialization.ParamsServices{
			Offsets: helpers.TimestampOffsets{
				OffsetUTCHours: +3,
			},
			APIKeyGeolocation: params.KeyGeolocationAPI,

			ServiceLogging: serviceLogging,
		},
	)
	if errInitialization != nil {
		fmt.Printf(
			"error create listener: %s\n",
			errInitialization.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConnectivityIssues,
		)
	}

	transportTCP, errCrTransport := transporttcp.NewTransportTCP(
		listener,
		&transporttcp.PiersNewTransportTCP{
			ServiceLogging:   serviceLogging,
			ServiceAnalytics: serviceAnalytics,
		},
	)
	if errCrTransport != nil {
		fmt.Printf(
			"error create transport TCP: %s\n",
			errCrTransport.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConnectivityIssues,
		)
	}

	return &App{
		transportHTTP: fiber.New(
			fiber.Config{
				BodyLimit: 1 * 1024 * 1024, // in mb
			},
		),

		transportTCP: transportTCP,

		serviceAnalytics: serviceAnalytics,
		serviceLogging:   serviceLogging,

		fnFreeResources: fnCloseLogging,

		portHTTP: params.ConfigPortHTTP,
	}
}
