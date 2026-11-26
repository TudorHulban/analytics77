package appanalytics

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (a *App) HandlerViewRegistry(c fiber.Ctx) error {
	c.Set("Content-Type", "text/html")

	aggs, errAgg := a.datacenter.CurrentMonthAggregateTopN("s2")
	if errAgg != nil {
		c.Status(http.StatusInternalServerError)

		return c.SendString(errAgg.Error())
	}

	data := []string{
		a.datacenter.String(),
	}

	for _, agg := range aggs {
		data = append(data, agg.String())
	}

	return c.SendString(
		strings.Join(data, "\n"),
	)
}
