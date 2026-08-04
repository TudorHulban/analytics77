package appanalytics

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (a *App) HandlerViewRegistry(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	data := []string{
		a.serviceAnalytics.DC.String(),
	}

	return c.SendString(
		strings.Join(data, "\n"),
	)
}
