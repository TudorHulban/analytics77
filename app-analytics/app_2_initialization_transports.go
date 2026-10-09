package appanalytics

import (
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func InitializeTransportHTTPRoutes(app *App) {
	app.transportHTTP.Get(
		"/ws",
		websocket.New(app.transportWS.HandleWebSocket),
	)

	app.transportHTTP.Get(
		_RoutesHTTP,
		app.handlerPage,
	)
}

func InitializeTransportWS(app *App) {
	// app.transportWS.Handlers["/login"] = app.wslogin
	app.transportWS.Handlers[WSRouteToday] = app.wsSearch
	app.transportWS.Handlers[WSRouteSync] = app.wsSearch

	app.transportHTTP.Use(
		"/ws",
		func(c fiber.Ctx) error {
			if c.Get("Upgrade") == "websocket" {
				return c.Next()
			}

			return fiber.ErrUpgradeRequired
		},
	)
}
