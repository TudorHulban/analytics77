package fragments

import (
	"log"
	"testing"

	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
)

func TestTopBar(t *testing.T) {
	el := TopBar{
		SiteCombo: SiteCombo{
			Status: "Active",
			Options: []inputs.Option{
				{
					Value: "site1",
					Label: "🌐 example.com",
				},
				{
					Value: "site2",
					Label: "🛒 shop.example.com",
				},
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
