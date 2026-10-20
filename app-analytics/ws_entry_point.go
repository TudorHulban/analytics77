package appanalytics

import (
	"fmt"

	"github.com/TudorHulban/hxgo/helpers/ws"
)

// func extractCredentials(v url.Values) (string, string) {
// 	raw := v.Encode()

// 	var username, password string

// 	// SplitSeq yields index and substring
// 	for pair := range strings.SplitSeq(raw, "&") {
// 		key, val, couldCut := strings.Cut(pair, "=")
// 		if !couldCut {
// 			continue
// 		}

// 		if key == "username" {
// 			username = val
// 		}

// 		if key == "password" {
// 			password = val
// 		}
// 	}

// 	return username, password
// }

// func (a *App) wslogin(c *ws.Client, message *ws.WSMessage) {
// 	user, password := extractCredentials(message.Values)

// 	if user == "admin" && password == "password" {
// 		_ = c.SendString(
// 			ws.WrapResponse(
// 				message.RequestID,
// 				fmt.Sprintf(
// 					`<div id="redirect" hx-redirect="%s"></div>`,
// 					_RouteAuthorised,
// 				),
// 			),
// 		)
// 	}

// 	_ = c.SendString(
// 		ws.WrapResponse(
// 			message.RequestID,
// 			fmt.Sprintf(
// 				`<div id="redirect" hx-redirect="%s"></div>`,
// 				_RouteNotAuthorised,
// 			),
// 		),
// 	)
// }

func (*App) wsLogger(c *ws.Client, message *ws.WSMessage) {
	fmt.Println("ws message: ", message.String())

	extractDateSelectors(message.Values)

	_ = c.SendString(
		ws.WrapResponse(
			message.RequestID,
			"",
		),
	)
}
