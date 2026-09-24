package fragments

import (
	"log"
	"testing"

	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
)

var testTopBar = func() TopBar {
	return TopBar{
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
			SelectedValue: "2026-4",
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

		SelectDay: ComboSelectDay{
			SelectedValue: "17",
			OptionsDay: []inputs.Option{
				{
					Value: "15",
					Label: "15",
				},
				{
					Value: "16",
					Label: "16",
				},
				{
					Value: "17",
					Label: "17",
				},
				{
					Value: "18",
					Label: "18",
				},
			},
		},

		SelectHour: ComboSelectHour{
			SelectedValue: "5",
			OptionsHour: []inputs.Option{
				{
					Value: "1",
					Label: "1",
				},
				{
					Value: "2",
					Label: "2",
				},
				{
					Value: "3",
					Label: "3",
				},
				{
					Value: "4",
					Label: "4",
				},
				{
					Value: "5",
					Label: "5",
				},
				{
					Value: "6",
					Label: "6",
				},
				{
					Value: "7",
					Label: "7",
				},
				{
					Value: "8",
					Label: "8",
				},
			},
		},
	}
}

func TestTopBar(t *testing.T) {
	el := testTopBar()

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
