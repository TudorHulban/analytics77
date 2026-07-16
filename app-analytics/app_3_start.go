package appanalytics

import (
	"fmt"
	"runtime"
)

func (a *App) root() string {
	return ":" + a.portHTTP
}

func (a *App) Start() error {
	messageStart := fmt.Sprintf(
		"starting application %s on %d threads",

		_NameApp,
		runtime.GOMAXPROCS(0),
	)

	a.serviceLogging.Logger.Info(messageStart)

	InitializeTransportRoutes(a)

	chError := make(chan error, 2)

	go func() {
		chError <- a.transportHTTP.Listen(a.root())
	}()

	go func() {
		chError <- a.transportTCP.Start()
	}()

	// Block and wait for the first error to arrive
	for range 2 {
		if errTransport := <-chError; errTransport != nil {
			// Log and return the first error that dropped
			return errTransport
		}
	}

	return nil
}
