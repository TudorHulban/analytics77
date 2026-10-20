package fragments

import (
	"strings"

	"github.com/TudorHulban/hxgo/dsl"
	"github.com/tudorhulban/hxhelpers"
)

type Metric struct {
	Name  string
	Value string
}

func (elem Metric) Build() dsl.Node {
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

type MetricCard struct {
	Title        string
	FAwesomeIcon string // ex. fa-network-wired
	CSSID        string

	Metrics []Metric
}

func (elem MetricCard) Build() dsl.Node {
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

		dsl.Div(
			dsl.AttrID(elem.CSSID),
			dsl.Ul(
				append([]dsl.Node{dsl.AttrClass(_CSSClassCardMetric)}, entries...)...,
			),
		),
	)
}

func (elem MetricCard) WSUpdateEntries() string {
	var sb strings.Builder

	sb.WriteString(
		hxhelpers.Sprintf(
			`<div id="%s">`,
			elem.CSSID),
	)

	sb.WriteString(
		hxhelpers.Sprintf(
			`<ul class="%s">`,
			_CSSClassCardMetric),
	)

	for _, metric := range elem.Metrics {
		sb.WriteString(
			hxhelpers.Sprintf(
				`<li class="top-item">`+
					`<span class="item-name">`+
					`<i class="fas fa-circle" style="font-size:0.5rem;color:var(--accent-primary);"></i>`+
					`%s`+
					`</span>`+
					`<span class="item-value">%s</span>`+
					`</li>`,

				metric.Name,
				metric.Value,
			),
		)
	}

	sb.WriteString(`</ul>`)
	sb.WriteString(`</div>`)

	return sb.String()
}

type MetricsContainer struct {
	CSSID string
	Cards []MetricCard
}

func (elem MetricsContainer) Build() dsl.Node {
	cards := make([]dsl.Node, len(elem.Cards))

	for ix, card := range elem.Cards {
		cards[ix] = card.Build()
	}

	return dsl.Div(
		append(
			[]dsl.Node{
				dsl.AttrID(elem.CSSID),
				dsl.AttrClass("metrics-grid"),
			}, cards...,
		)...,
	)
}
