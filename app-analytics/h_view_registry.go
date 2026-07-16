package appanalytics

import "github.com/gofiber/fiber/v3"

func (a *App) HandlerViewRegistry(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	return c.SendString(
		a.serviceAnalytics.DC.String(),
	)
}
