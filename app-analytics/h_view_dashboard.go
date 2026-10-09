package appanalytics

import (
	"time"

	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/templates/fragments"
	"github.com/tudorhulban/analytics77/templates/layouts"
)

func (*App) handlerPage(c fiber.Ctx) error {
	now := time.Now()

	c.Set("Content-Type", "text/html")

	top := fragments.TopBar{
		URLBtnToday: WSRouteToday,
		URLBtnSync:  WSRouteSync,

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

		SelectMonth: fragments.NewComboSelectMonth(now, 3, 2),
		SelectDay:   fragments.NewComboSelectDay(now, 7, 2),
		SelectHour:  fragments.NewComboSelectHour(now, 14, 3),

		CounterRecords: fragments.Counter{
			CSSDivClass: "records-counter",
			SpanID:      _CSSIdCounterRecordsValue,

			FAwesomeIcon: "fa-database",
			Label:        "1777",
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
