package layouts

import "github.com/TudorHulban/hxgo/dsl"

var _DivApp = dsl.Div(
	dsl.AttrID("appScreen"),
	_HeaderMobile,
)

var _HeaderMobile = dsl.Header(
	dsl.Class("mobile-header"),
	dsl.Div(
		dsl.Class("logo-mobile"),
		dsl.I(
			dsl.Class("fas fa-shield-haltered"),
		),
		dsl.Span(
			dsl.Text("MetricActive"),
		),
	),
	dsl.Div(
		dsl.Class("mobile-header-actions"),
		dsl.Button(
			dsl.Class("icon-btn"),
			dsl.AttrID("mobileThemeToggle"),
			dsl.Title(
				dsl.Text("Toggle theme"),
			),
			dsl.I(
				dsl.Class("fas fa-moon"),
			),
		),
	),
)
