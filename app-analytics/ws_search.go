package appanalytics

import (
	"fmt"

	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/tudorhulban/analytics77/templates/fragments"
)

func (a *App) wsSearch(c *ws.Client, message *ws.WSMessage) {
	fmt.Println("ws message: ", message.String())

	extractDateSelectors(message.Values)

	_ = c.SendString(
		ws.WrapResponse(
			message.RequestID,
			fragments.Counter{
				SpanID: _CSSIdCounterRecordsValue,
				Label:  "77",
			}.
				RawString(),
		),
	)
}
