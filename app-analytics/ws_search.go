package appanalytics

import (
	"time"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/TudorHulban/hxgo/helpers/ws"
	"github.com/tudorhulban/analytics77/templates/fragments"
)

// fmt.Println("ws message: ", message.String())
// extractDateSelectors(message.Values)

func (*App) wsSearch(c *ws.Client, message *ws.WSMessage) {
	_ = c.SendString(
		ws.WrapResponses(
			message.RequestID,
			fragments.Counter{
				SpanID: _CSSIdCounterRecordsValue,
				Label:  "77",
			}.
				WSUpdateLabel(),
		),
	)
}

// fmt.Println("ws message: ", message.String())
// extractDateSelectors(message.Values)

func (*App) wsToday(c *ws.Client, message *ws.WSMessage) {
	_ = c.Send(
		ws.WrapByteResponses(
			message.RequestID,
			[]byte(
				fragments.Counter{
					SpanID: _CSSIdCounterRecordsValue,
					Label:  "77",
				}.
					WSUpdateLabel(),
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

			[]byte(
				fragments.MetricCard{
					CSSID: _CSSIDCard01,

					Metrics: []fragments.Metric{
						{
							Name:  "Value 1",
							Value: randomNumberString(1, 4),
						},
						{
							Name:  "Value 2",
							Value: randomNumberString(1, 4),
						},
					},
				}.
					WSUpdateEntries(),
			),

			[]byte(
				fragments.MetricCard{
					CSSID: _CSSIDCard02,

					Metrics: []fragments.Metric{
						{
							Name:  "Value A",
							Value: randomNumberString(1, 4),
						},
						{
							Name:  "Value B",
							Value: randomNumberString(1, 3),
						},
					},
				}.
					WSUpdateEntries(),
			),
		),
	)
}
