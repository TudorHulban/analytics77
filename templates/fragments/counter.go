package fragments

import (
	"github.com/TudorHulban/hxgo/dsl"
)

type Counter struct {
	CSSDivID    string
	CSSDivClass string
	SpanID      string // needed in RawSelect.

	FAwesomeIcon string
	Label        string // updated through RawSelect.
}

func (elem Counter) Build() dsl.Node {
	return dsl.Div(
		dsl.AttrClass(elem.CSSDivClass),
		dsl.AttrIDIf(
			len(elem.CSSDivID) > 0,
			elem.CSSDivID,
		),

		dsl.I(
			dsl.AttrClass("fas "+elem.FAwesomeIcon),
		),
		dsl.Span(
			dsl.AttrIDIf(
				len(elem.SpanID) > 0,
				elem.SpanID,
			),
			dsl.Text(elem.Label),
		),
	)
}

// RawSelect knows for the element how to construct the update.
// It returns the inner HTML for hx match.
//
// The element might be even zero besides the updated value.
func (elem Counter) RawSelect() dsl.Node {
	return dsl.Span(
		dsl.AttrIDIf(
			len(elem.SpanID) > 0,
			elem.SpanID,
		),
		dsl.Text(elem.Label),
	)
}
