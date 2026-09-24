package fragments

import (
	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
)

type ComboSelectMonth struct {
	SelectedValue string
	OptionsMonth  []inputs.Option
}

func (elem *ComboSelectMonth) Build() dsl.Node {
	if len(elem.OptionsMonth) == 0 {
		return dsl.Node{}
	}

	return dsl.Div(
		dsl.AttrClass("selector-group"),
		dsl.AttrID("monthGroup"),

		inputs.InputSelect{
			CSSDivID: "monthSelect",

			SelectedValue: elem.SelectedValue,
			SelectOptions: elem.OptionsMonth,
		}.
			Raw(),
	)
}

type ComboSelectDay struct {
	SelectedValue string
	OptionsDay    []inputs.Option
}

func (elem *ComboSelectDay) Build() dsl.Node {
	if len(elem.OptionsDay) == 0 {
		return dsl.Node{}
	}

	return dsl.Div(
		dsl.AttrClass("selector-group"),
		dsl.AttrID("monthGroup"),

		inputs.InputSelect{
			CSSDivID: "daySelect",

			SelectedValue: elem.SelectedValue,
			SelectOptions: elem.OptionsDay,
		}.
			Raw(),
	)
}

type ComboSelectHour struct {
	SelectedValue string
	OptionsHour   []inputs.Option
}

func (elem *ComboSelectHour) Build() dsl.Node {
	if len(elem.OptionsHour) == 0 {
		return dsl.Node{}
	}

	return dsl.Div(
		dsl.AttrClass("selector-group"),
		dsl.AttrID("monthGroup"),

		inputs.InputSelect{
			CSSDivID: "hourSelect",

			SelectedValue: elem.SelectedValue,
			SelectOptions: elem.OptionsHour,
		}.
			Raw(),
	)
}

type ComboSite struct {
	Status      string
	OptionsSite []inputs.Option
}

func (elem *ComboSite) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass("active-site-combo"),
		dsl.I(
			dsl.AttrClass("fas fa-globe"),
		),

		inputs.InputSelect{
			CSSDivID: "siteSelector",

			SelectOptions: elem.OptionsSite,
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
					"month",
				),
				dsl.I(
					dsl.AttrClass("far fa-calendar-alt"),
				),
			),
		),
	)
}

type TopBar struct {
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
			),
		),
	)
}
