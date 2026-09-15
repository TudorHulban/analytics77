package appanalytics

import (
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/templates/layouts"
)

func (*App) handlerPage(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	h := dsl.Label(
		dsl.Text("xxxxxxxxxxxxxxx"),
	)

	p := layouts.Page{
		Title:       "page title",
		Description: "page description",
		Language:    _PageLanguageEnglish,

		Body: []dsl.Node{
			layouts.LayoutMobile(h),
		},
	}

	content := dsl.RenderFast(p.Build())

	// fmt.Println(string(content))

	return c.Send(
		content,
	)
}
