package layouts

import (
	"log"
	"testing"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
)

func TestBase(t *testing.T) {
	el := Page{
		Title: "Title",

		Body: []dsl.Node{
			LayoutMobile(dsl.Node{}),
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
