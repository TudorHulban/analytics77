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
		ws.WrapResponses(
			message.RequestID,
			fragments.Counter{
				SpanID: _CSSIdCounterRecordsValue,
				Label:  "77",
			}.
				ResponseUpdateLabel(),
		),
	)
}

func (a *App) wsToday(c *ws.Client, message *ws.WSMessage) {
	fmt.Println("ws message: ", message.String())

	extractDateSelectors(message.Values)

	_ = c.SendString(
		ws.WrapResponses(
			message.RequestID,
			fragments.Counter{
				SpanID: _CSSIdCounterRecordsValue,
				Label:  "77",
			}.
				ResponseUpdateLabel(),
		),
	)
}
