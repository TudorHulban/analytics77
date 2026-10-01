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
		dsl.AttrID(_IDComboSelectMonth),

		inputs.InputSelect{
			CSSInputID: _IDSelectMonth,

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
		dsl.AttrID(_IDComboSelectDay),

		inputs.InputSelect{
			CSSInputID: _IDSelectDay,

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
		dsl.AttrID(_IDComboSelectHour),

		inputs.InputSelect{
			CSSInputID: _IDSelectHour,

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
