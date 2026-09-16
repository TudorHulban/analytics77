package fragments

import (
	"log"
	"testing"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
)

func TestCardMetric(t *testing.T) {
	el := CardMetric{
		Title:        "Title",
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
