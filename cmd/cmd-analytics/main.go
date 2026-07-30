package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	// Context that cancels on SIGINT or SIGTERM
	ctx, stop := signal.NotifyContext(
		context.Background(),

		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// If a SECOND signal arrives while we're still shutting down,
	// bail out immediately instead of hanging forever.
	go func() {
		<-ctx.Done() // first Ctrl+C: begin graceful shutdown
		stop()       // restore default signal behavior

		chSignal := make(chan os.Signal, 1)

		signal.Notify(chSignal, os.Interrupt, syscall.SIGTERM)
		<-chSignal // second Ctrl+C: force it

		fmt.Println("\nforced shutdown, exiting immediately")

		os.Exit(1)
	}()

	// context.Canceled is not an error — it is the shutdown request itself
	if errStart := app.Start(ctx); errStart != nil {
		fmt.Printf(
			"error application start: %s\n",
			errStart.Error(),
		)

		os.Exit(1)
	}
}
