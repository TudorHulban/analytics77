package main

import (
	"fmt"
	"os"

	appanalytics "github.com/tudorhulban/analytics77/app-analytics"
	"github.com/tudorhulban/analytics77/cmd"
	"github.com/tudorhulban/analytics77/infra/initialization"
	"github.com/tudorhulban/hxerrors"
)

func main() {
	configRaw := initialization.Configuration(cmd.PathConfig)

	configuration, errParse := extractConfiguration(configRaw)
	if errParse != nil {
		fmt.Printf(
			"error extract configuration: %s\n",
			errParse.Error(),
		)

		os.Exit(
			hxerrors.OSExitForConfigurationIssues,
		)
	}

	app := appanalytics.InitializeApp(
		&appanalytics.ParamsInitializeApp{
			ConfigPortRPC:  configuration.portRPC,
			ConfigPortHTTP: configuration.portHTTP,

			PathLogFile:       configuration.nameLogfile,
			KeyGeolocationAPI: os.Getenv(cmd.OSAPIGeolocation),
		},

		&appanalytics.PiersInitializeApp{
			Writer:   os.Stderr,
			FuncExit: os.Exit,
		},
	)

	fmt.Println(
		app.Start(),
	)

	// TODO: add gracefully shutdown support
}
