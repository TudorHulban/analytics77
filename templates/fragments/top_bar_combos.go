package fragments

import (
	"fmt"
	"time"

	"github.com/TudorHulban/hxgo/components/inputs"
	"github.com/TudorHulban/hxgo/dsl"
)

type ComboSelectMonth struct {
	SelectedValue string
	OptionsMonth  []inputs.Option
}

func NewComboSelectMonth(t time.Time, monthsBehind, monthsAhead int) ComboSelectMonth {
	selectedValue := fmt.Sprintf("%d-%d", t.Year(), t.Month()) // Format as ("YYYY-M").

	var options []inputs.Option

	// Loop from negative (behind) to positive (ahead)
	for ix := -monthsBehind; ix <= monthsAhead; ix++ {
		targetTime := t.AddDate(0, ix, 0)

		options = append(
			options,
			inputs.Option{
				Value: fmt.Sprintf(
					"%d-%d",
					targetTime.Year(), targetTime.Month(),
				),
				Label: targetTime.Format("Jan 2006"), // e.g., "Apr 2026", "May 2026"
			},
		)
	}

	return ComboSelectMonth{
		SelectedValue: selectedValue,
		OptionsMonth:  options,
	}
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

func NewComboSelectDay(t time.Time, daysBehind, daysAhead int) ComboSelectDay {
	selectedValue := t.Format("2006-01-02") // Full date to match option values

	var options []inputs.Option

	// Loop from negative (behind) to positive (ahead)
	for ix := -daysBehind; ix <= daysAhead; ix++ {
		targetTime := t.AddDate(0, 0, ix)

		options = append(
			options,
			inputs.Option{
				Value: targetTime.Format("2006-01-02"), // e.g., "2026-09-30", "2026-10-04"
				Label: targetTime.Format("02"),         // e.g., "30", "04" (or targetTime.Format("Jan 02") for month context)
			},
		)
	}

	return ComboSelectDay{
		SelectedValue: selectedValue,
		OptionsDay:    options,
	}
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

func NewComboSelectHour(t time.Time, hoursBehind, hoursAhead int) ComboSelectHour {
	selectedValue := t.Format("15") // Format as 24-hour (e.g., "14")

	var options []inputs.Option

	// Loop from negative (behind) to positive (ahead)
	for ix := -hoursBehind; ix <= hoursAhead; ix++ {
		targetTime := t.Add(time.Duration(ix) * time.Hour)

		options = append(
			options,
			inputs.Option{
				Value: targetTime.Format("15"),
				Label: targetTime.Format("15:00"), // e.g., "14:00", "15:00"
			},
		)
	}

	return ComboSelectHour{
		SelectedValue: selectedValue,
		OptionsHour:   options,
	}
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
