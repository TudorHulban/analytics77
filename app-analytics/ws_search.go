package appanalytics

import (
	"fmt"

	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/tudorhulban/analytics77/templates/fragments"
)

func (a *App) wsSearch(c *websocket.Conn, message *ws.WSMessage) {
	fmt.Println("ws message: ", message.String())

	extractDateSelectors(message.Values)

	counter := fragments.Counter{
		SpanID: _CSSIdCounterRecordsValue,
		Label:  "77",
	}

	_ = c.WriteMessage(
		websocket.TextMessage,
		[]byte(
			a.transportWS.WrapResponse(
				message.RequestID,
				"",
			),
		),
	)
}
