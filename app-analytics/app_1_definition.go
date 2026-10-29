package appanalytics

import (
	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/infra/datacenter"
	transporttcp "github.com/tudorhulban/analytics77/infra/transport-tcp"
	"github.com/tudorhulban/analytics77/services/sgeo"
	"github.com/tudorhulban/analytics77/services/slogging"
)

type App struct {
	transportHTTP *fiber.App
	transportWS   *ws.ServerWS
	transportTCP  *transporttcp.TransportTCP

	serviceLogging *slogging.ServiceLogging
	serviceGeo     *sgeo.ServiceGeo

	fnFreeResources func()

	datacenter *datacenter.DataCenter

	portHTTP string
}
