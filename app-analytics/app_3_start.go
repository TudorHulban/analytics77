package appanalytics

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"sync"
)

func (a *App) root() string {
	return ":" + a.portHTTP
}

func (a *App) Start(ctx context.Context) error {
	messageStart := fmt.Sprintf(
		"starting application %s on %d threads",

		_NameApp,
		runtime.GOMAXPROCS(0),
	)

	a.serviceLogging.Logger.Info(messageStart)

	InitializeTransportRoutes(a)

	chError := make(chan error, 2)

	var wgStartElements sync.WaitGroup
	wgStartElements.Add(2)

	startElement := func(start func() error) {
		defer wgStartElements.Done()

		if errStart := start(); errStart != nil {
			chError <- errStart
		}
	}

	go startElement(
		func() error {
			errTransport := a.transportHTTP.Listen(a.root())
			if errors.Is(errTransport, http.ErrServerClosed) {
				return nil
			}

			return errTransport
		},
	)

	go startElement(
		func() error {
			return a.transportTCP.Start()
		},
	)

	var stopOnce sync.Once

	stopElements := func() {
		a.transportTCP.Stop()
		a.Stop()
	}

	context.AfterFunc(
		ctx,
		func() {
			stopOnce.Do(stopElements)
		},
	)

	var errStartElement error

	select {
	case errStartElement = <-chError:

	case <-ctx.Done():
		a.serviceLogging.Logger.Print(
			"shutdown in progress based on context cancellation",
		)

		stopOnce.Do(stopElements)
		wgStartElements.Wait()

		return nil // graceful shutdown is not an error
	}

	stopOnce.Do(stopElements)

	wgStartElements.Wait()

	return errStartElement
}
