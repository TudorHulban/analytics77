package appanalytics

import (
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

	return c.Send(nil)
}
