package fragments

import (
	"log"
	"testing"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
)

var (
	testCard1 = func() MetricCard {
		return MetricCard{
			Title:        "Card 1",
			FAwesomeIcon: "fa-network-wired",

			Metrics: []Metric{
				{
					Name:  "Product A",
					Value: "21",
				},
				{
					Name:  "Product B",
					Value: "11",
				},
			},
		}
	}

	testCard2 = func() MetricCard {
		return MetricCard{
			Title:        "Card 2",
			FAwesomeIcon: "fa-network-wired",

			Metrics: []Metric{
				{
					Name:  "Product C",
					Value: "71",
				},
				{
					Name:  "Product D",
					Value: "31",
				},
			},
		}
	}
)

func TestCardMetric(t *testing.T) {
	el := testCard1()

	app := fiber.New()

	app.Get(
		"/",
		func(c fiber.Ctx) error {
			c.Type("html")

			return c.Send(
				dsl.RenderFast(el.Build()),
			)
		},
	)

	log.Fatal(app.Listen(":3000"))
}

func Test_Container_CardMetric(t *testing.T) {
	el := MetricsContainer{
		CSSID: "metricsContainer",
		Cards: []MetricCard{
			testCard1(),
			testCard2(),
		},
	}

	app := fiber.New()

	app.Get(
		"/",
		func(c fiber.Ctx) error {
			c.Type("html")

			return c.Send(
				dsl.RenderFast(el.Build()),
			)
		},
	)

	log.Fatal(app.Listen(":3000"))
}
