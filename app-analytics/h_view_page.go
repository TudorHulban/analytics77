package appanalytics

import (
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/templates/fragments"
	"github.com/tudorhulban/analytics77/templates/layouts"
)

func (*App) handlerPage(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

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
			layouts.LayoutMobile(h.Build()),
		},
	}

	content := dsl.RenderFast(p.Build())

	// fmt.Println(string(content))

	return c.Send(
		content,
	)
}
