package appanalytics

import (
	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/gofiber/fiber/v3"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/slogging"
)

type App struct {
	transportHTTP *fiber.App
	transportWS   *ws.ServerWS
	transportTCP  *transporttcp.TransportTCP

	serviceLogging   *slogging.ServiceLogging
	serviceAnalytics *sanalytics.ServiceAnalytics

	fnFreeResources func()

	portHTTP string
}
