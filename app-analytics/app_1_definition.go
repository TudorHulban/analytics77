package appanalytics

import (
	"github.com/gofiber/fiber/v3"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/sanalytics"
	"github.com/tudorhulban/analytics77/services/slogging"
)

type App struct {
	transportHTTP *fiber.App
	transportTCP  *transporttcp.TransportTCP

	serviceLogging   *slogging.ServiceLogging
	serviceAnalytics *sanalytics.ServiceAnalytics

	fnFreeResources func()

	portHTTP string
}
