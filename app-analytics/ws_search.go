package appanalytics

import (
	"fmt"
	"time"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/tudorhulban/analytics77/templates/fragments"
)

func (*App) wsSearch(c *ws.Client, message *ws.WSMessage) {
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

func (*App) wsToday(c *ws.Client, message *ws.WSMessage) {
	// fmt.Println("ws message: ", message.String())

	// extractDateSelectors(message.Values)

	_ = c.Send(
		ws.WrapByteResponses(
			message.RequestID,
			[]byte(
				fragments.Counter{
					SpanID: _CSSIdCounterRecordsValue,
					Label:  "77",
				}.
					ResponseUpdateLabel(),
			),

			dsl.RenderHX(
				fragments.NewComboSelectMonth(
					time.Now(),
					_DefaultMonthsBehind,
					_DefaultMonthsAhead,
				).
					Build(),

				fragments.NewComboSelectDay(
					time.Now(),
					_DefaultDaysBehind,
					_DefaultDaysAhead,
				).
					Build(),

				fragments.NewComboSelectHour(
					time.Now(),
					_DefaultDaysBehind,
					_DefaultDaysAhead,
				).
					Build(),
			),
		),
	)
}
