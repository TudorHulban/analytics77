package layouts

import "github.com/TudorHulban/hxgo/dsl"

var _MobileHeader = dsl.Header(
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
		dsl.Button(
			dsl.Class("icon-btn"),
			dsl.Title(
				dsl.Text("Logout"),
			),
			dsl.AttrHXRedirect("/logout"),
			dsl.I(
				dsl.Class("fas fa-sign-out-alt"),
			),
		),
		dsl.Button(
			dsl.Class("icon-btn"),
			dsl.AttrID("hamburgerBtn"),
			dsl.Title(
				dsl.Text("Menu"),
			),
		),
	),
)

var _MobileBottomNavigation = dsl.Nav(
	dsl.Class("mobile-bottom-nav"),
	dsl.Div(
		dsl.Class("mobile-nav-items"),
		dsl.Div(
			dsl.Class("mobile-nav-item active"),
			dsl.AttrHXGET("/analytics"),
			dsl.I(
				dsl.Class("fas fa-chart-line"),
			),
			dsl.Span(
				dsl.Text("Analytics"),
			),
		),
		dsl.Div(
			dsl.Class("mobile-nav-item"),
			dsl.AttrHXGET("/users"),
			dsl.I(
				dsl.Class("fas fa-users"),
			),
			dsl.Span(
				dsl.Text("Users"),
			),
		),
	),
)
