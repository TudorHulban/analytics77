package fragments

import (
	"github.com/TudorHulban/hxgo/dsl"
)

type PeriodToggle struct {
	ShowHour  bool
	ShowDay   bool
	ShowMonth bool
}

func (elem *PeriodToggle) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass("period-toggle"),

		dsl.If(
			elem.ShowHour,

			dsl.Button(
				dsl.AttrClass("period-btn active"),
				dsl.AttrWithValue(
					"data-period",
					"hour",
				),
				dsl.I(
					dsl.AttrClass("far fa-clock"),
				),
				dsl.Text("Hour"),
			),
		),

		dsl.If(
			elem.ShowDay,

			dsl.Button(
				dsl.AttrClass("period-btn"),
				dsl.AttrWithValue(
					"data-period",
					"day",
				),
				dsl.I(
					dsl.AttrClass("far fa-calendar-check"),
				),
				dsl.Text("Day"),
			),
		),

		dsl.If(
			elem.ShowMonth,

			dsl.Button(
				dsl.AttrClass("period-btn"),
				dsl.AttrWithValue(
					"data-period",
					"month",
				),
				dsl.I(
					dsl.AttrClass("far fa-calendar-alt"),
				),
				dsl.Text("Month"),
			),
		),
	)
}

type TopBar struct {
	CounterRecords Counter

	SiteCombo ComboSite

	SelectMonth ComboSelectMonth
	SelectDay   ComboSelectDay
	SelectHour  ComboSelectHour

	PeriodToggler PeriodToggle
}

func (elem *TopBar) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass("top-bar"),

		elem.SiteCombo.Build(),

		dsl.Div(
			dsl.AttrClass("controls-wrapper"),

			elem.PeriodToggler.Build(),

			dsl.Div(
				dsl.AttrClass("date-selectors"),

				dsl.If(
					elem.PeriodToggler.ShowMonth,
					elem.SelectMonth.Build(),
				),

				dsl.If(
					elem.PeriodToggler.ShowDay,
					elem.SelectDay.Build(),
				),

				dsl.If(
					elem.PeriodToggler.ShowHour,
					elem.SelectHour.Build(),
				),

				dsl.Button(
					dsl.AttrClass("quick-day-btn"),
					dsl.AttrID("quickDayBtn"),

					dsl.I(
						dsl.AttrClass("fas fa-bolt"),
					),

					dsl.Text("Today"),
				),

				elem.CounterRecords.Build(),
			),
		),
	)
}
