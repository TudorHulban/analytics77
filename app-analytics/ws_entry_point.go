package appanalytics

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/gofiber/contrib/v3/websocket"
)

func extractCredentials(v url.Values) (string, string) {
	raw := v.Encode()

	var username, password string

	// SplitSeq yields index and substring
	for pair := range strings.SplitSeq(raw, "&") {
		key, val, couldCut := strings.Cut(pair, "=")
		if !couldCut {
			continue
		}

		if key == "username" {
			username = val
		}

		if key == "password" {
			password = val
		}
	}

	return username, password
}

func (a *App) wslogin(c *websocket.Conn, message *ws.WSMessage) {
	user, password := extractCredentials(message.Values)

	if user == "admin" && password == "password" {
		_ = c.WriteMessage(
			websocket.TextMessage,
			[]byte(
				a.transportWS.WrapResponse(
					message.RequestID,
					fmt.Sprintf(
						`<div id="redirect" hx-redirect="%s"></div>`,
						_RouteAuthorised,
					),
				),
			),
		)
	}

	_ = c.WriteMessage(
		websocket.TextMessage,
		[]byte(
			a.transportWS.WrapResponse(
				message.RequestID,
				fmt.Sprintf(
					`<div id="redirect" hx-redirect="%s"></div>`,
					_RouteNotAuthorised,
				),
			),
		),
	)
}
