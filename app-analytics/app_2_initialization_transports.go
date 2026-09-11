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
		_RoutesAll,
		app.handlerPage,
	)
}

func InitializeTransportWS(app *App) {
	app.transportWS.Handlers["/login"] = app.wslogin

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
