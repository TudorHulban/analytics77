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
	if len(os.Args) < 2 {
		fmt.Println(
			"Error: Please provide the geolocation API key as the first argument.",
		)
		fmt.Println(
			"Usage: go run main.go <API_KEY>",
		)

		os.Exit(
			hxerrors.OSExitForApplicationIssues,
		)
	}

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
			KeyGeolocationAPI: os.Args[1],
		},
	)

	fmt.Println(
		app.Start(),
	)

	// TODO: add gracefully shutdown support
}
