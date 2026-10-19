package appanalytics

// func (*App) handlerSearch(c fiber.Ctx) error {
// 	defer hxhelpers.TraceExitTime(time.Now())

// 	c.Set("Content-Type", "text/html")

// 	return c.Send(
// 		dsl.RenderHX(
// 			fragments.Counter{
// 				SpanID: _CSSIdCounterRecordsValue,
// 				Label:  "77",
// 			}.
// 				HTTPUpdateLabel(),

// 			fragments.NewComboSelectMonth(
// 				time.Now(),
// 				_DefaultMonthsBehind,
// 				_DefaultMonthsAhead,
// 			).
// 				Build(),

// 			fragments.NewComboSelectDay(
// 				time.Now(),
// 				_DefaultDaysBehind,
// 				_DefaultDaysAhead,
// 			).
// 				Build(),

// 			fragments.NewComboSelectHour(
// 				time.Now(),
// 				_DefaultDaysBehind,
// 				_DefaultDaysAhead,
// 			).
// 				Build(),
// 		),
// 	)
// }
