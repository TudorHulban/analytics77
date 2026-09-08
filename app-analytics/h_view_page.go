package appanalytics

import (
	"github.com/TudorHulban/hxgo/dsl"
	"github.com/gofiber/fiber/v3"
	"github.com/tudorhulban/analytics77/templates/layouts"
)

func (a *App) HandlerPage(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	p := layouts.Page{
		Title:       "page title",
		Description: "page description",
		Language:    "Eng",
	}

	content := dsl.RenderFast(p.Build())

	// fmt.Println(string(content))

	return c.Send(
		content,
	)
}
