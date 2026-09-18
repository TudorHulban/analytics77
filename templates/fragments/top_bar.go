package fragments

import (
	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
)

type SiteCombo struct {
	Status  string
	Options []inputs.Option
}

func (elem *SiteCombo) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass("active-site-combo"),
		dsl.I(
			dsl.AttrClass("fas fa-globe"),
		),

		inputs.InputSelect{
			CSSDivID: "siteSelector",

			SelectOptions: elem.Options,
		}.
			Raw(),

		dsl.Div(
			dsl.AttrClass("site-status"),
			dsl.Span(
				dsl.AttrClass("status-dot"),
			),
			dsl.Span(
				dsl.Text(elem.Status),
			),
		),
	)
}

type PeriodToggle struct {
	ShowMonth bool
	ShowDay   bool
	ShowHour  bool
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
			),
		),

		dsl.If(
			elem.ShowMonth,

			dsl.Button(
				dsl.AttrClass("period-btn"),
				dsl.AttrWithValue(
					"data-period",
					"day",
				),
			),
		),
	)
}

type TopBar struct {
	SiteCombo     SiteCombo
	PeriodToggler PeriodToggle
}

func (elem *TopBar) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass("top-bar"),

		elem.SiteCombo.Build(),
	)
}
