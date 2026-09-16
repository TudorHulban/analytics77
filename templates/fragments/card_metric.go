package fragments

import "github.com/TudorHulban/hxgo/dsl"

type Metric struct {
	Name  string
	Value string
}

func (elem *Metric) Build() dsl.Node {
	return dsl.Li(
		dsl.AttrClass("top-item"),

		dsl.Span(
			dsl.AttrClass("item-name"),

			dsl.I(
				dsl.AttrClass("fas fa-circle"),
				dsl.AttrCSS("font-size:0.5rem;color:var(--accent-primary);"),
			),

			dsl.Text(elem.Name),
		),

		dsl.Span(
			dsl.AttrClass("item-value"),
			dsl.Text(elem.Value),
		),
	)
}

type CardMetric struct {
	Title        string
	FAwesomeIcon string // ex. fa-network-wired

	Metrics []Metric
}

func (elem *CardMetric) Build() dsl.Node {
	entries := make([]dsl.Node, len(elem.Metrics))

	for ix, metric := range elem.Metrics {
		entries[ix] = metric.Build()
	}

	return dsl.Div(
		dsl.AttrClass("metric-card"),

		dsl.Div(
			dsl.AttrClass("card-header"),
			dsl.Span(
				dsl.Text(elem.Title),
			),
			dsl.I(
				dsl.AttrClass("fas "+elem.FAwesomeIcon),
			),
		),

		dsl.Ul(
			append([]dsl.Node{dsl.AttrClass("top-list")}, entries...)...,
		),
	)
}
