package appanalytics

import (
	"fmt"
	"io"
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

	OffsetUTCHours int64
}

type PiersInitializeApp struct {
	Writer   io.Writer
	FuncExit func(int)
}

func InitializeApp(params *ParamsInitializeApp, piers *PiersInitializeApp) *App {
	if piers.Writer == nil {
		piers.Writer = os.Stderr
	}

	if piers.FuncExit == nil {
		piers.FuncExit = os.Exit
	}

	listener, errListener := net.Listen( //nolint:noctx
		"tcp",
		fmt.Sprintf(
			"127.0.0.1:%s",
			params.ConfigPortRPC,
		),
	)
	if errListener != nil {
		fmt.Fprintf(
			piers.Writer,
			"error create listener: %s\n",
			errListener.Error(),
		)

		piers.FuncExit(
			hxerrors.OSExitForConnectivityIssues,
		)
	}

	serviceLogging, fnCloseLogging, erCrServiceLogging := slogging.NewServiceLogging(params.PathLogFile, os.Stdout)
	if erCrServiceLogging != nil {
		fmt.Fprintf(
			piers.Writer,
			"error create servce logging: %s\n",
			erCrServiceLogging.Error(),
		)

		piers.FuncExit(
			hxerrors.OSExitForLoggingIssues,
		)
	}

	serviceAnalytics, errInitialization := initialization.Services(
		&initialization.ParamsServices{
			Offsets: helpers.TimestampOffsets{
				OffsetUTCHours: params.OffsetUTCHours,
			},
			APIKeyGeolocation: params.KeyGeolocationAPI,

			ServiceLogging: serviceLogging,
		},
	)
	if errInitialization != nil {
		fmt.Fprintf(
			piers.Writer,
			"error create listener: %s\n",
			errInitialization.Error(),
		)

		piers.FuncExit(
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
		fmt.Fprintf(
			piers.Writer,
			"error create transport TCP: %s\n",
			errCrTransport.Error(),
		)

		piers.FuncExit(
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
