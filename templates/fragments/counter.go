package fragments

import "github.com/TudorHulban/hxgo/dsl"

type Counter struct {
	CSSDivID    string
	CSSDivClass string
	SpanID      string

	FAwesomeIcon string
	Label        string
}

func (elem *Counter) Build() dsl.Node {
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
