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
		SiteCombo: ComboSite{
			Status: "Active",
			OptionsSite: []inputs.Option{
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

		PeriodToggler: PeriodToggle{
			ShowMonth: true,
			ShowDay:   true,
			ShowHour:  true,
		},

		SelectMonth: ComboSelectMonth{
			OptionsMonth: []inputs.Option{
				{
					Value: "2026-3",
					Label: "Apr 2026",
				},
				{
					Value: "2026-4",
					Label: "May 2026",
				},
				{
					Value: "2026-5",
					Label: "Jun 2026",
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
