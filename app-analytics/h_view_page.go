package appanalytics

import (
	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/templates/fragments"
	"github.com/tudorhulban/analytics77/templates/layouts"
)

func (*App) handlerPage(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	top := fragments.TopBar{
		SiteCombo: fragments.ComboSite{
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

		PeriodToggler: fragments.PeriodToggle{
			ShowMonth: true,
			ShowDay:   true,
			ShowHour:  true,
		},

		SelectMonth: fragments.ComboSelectMonth{
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

		SelectDay: fragments.ComboSelectDay{
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

		SelectHour: fragments.ComboSelectHour{
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

	h := fragments.MetricsContainer{
		CSSID: "metricsContainer",
		Cards: []fragments.MetricCard{
			{
				Title:        "Card 1",
				FAwesomeIcon: "fa-globe",

				Metrics: []fragments.Metric{
					{
						Name:  "Value 1",
						Value: "55",
					},
					{
						Name:  "Value 2",
						Value: "35",
					},
				},
			},
			{
				Title:        "Card 2",
				FAwesomeIcon: "fa-city",

				Metrics: []fragments.Metric{
					{
						Name:  "Value A",
						Value: "15",
					},
					{
						Name:  "Value B",
						Value: "3",
					},
				},
			},
		},
	}

	p := layouts.Page{
		Title:       "Dashboard",
		Description: "page description",
		Language:    _PageLanguageEnglish,

		Body: []dsl.Node{
			layouts.LayoutMobile(
				dsl.Div(
					dsl.AttrClass("page active"),
					dsl.AttrID("analytics-page"),

					top.Build(),
					h.Build(),
				),
			),
		},
	}

	content := dsl.RenderFast(p.Build())

	// fmt.Println(string(content))

	return c.Send(
		content,
	)
}
