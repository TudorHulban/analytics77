package appanalytics

import "fmt"

func (a *App) Stop() {
	a.fnFreeResources()

	a.serviceLogging.Logger.Print(
		"shutting down HTTP transport",
	)

	a.transportTCP.Stop()

	errStopFiber := a.transportHTTP.Shutdown()
	if errStopFiber != nil {
		a.serviceLogging.Logger.
			Error(
				fmt.Sprintf(
					"error shutting down HTTP transport: %v",
					errStopFiber,
				),
			)
	}
}
