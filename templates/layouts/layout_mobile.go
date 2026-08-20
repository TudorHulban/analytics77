package layouts

import "github.com/TudorHulban/hxgo/dsl"

func LayoutMobile(content dsl.Node) dsl.Node {
	return dsl.Div(
		dsl.AttrID("appScreen"),
		_MobileHeader,

		dsl.Div(
			dsl.Class("content-wrapper"),

			// sidebar
			dsl.Aside(
				dsl.Class("sidebar"),
				dsl.AttrID("sidebar"),

				dsl.Div(
					dsl.Class("logo"),
					dsl.I(
						dsl.Class("fas fa-shield-haltered"),
						dsl.Span(
							dsl.Text("MetricActive"),
						),
					),
				),

				// navigation links
				dsl.Nav(
					dsl.Class("nav-menu"),
					dsl.A(
						dsl.Class("nav-item active"),
						dsl.AttrHXGET("/analytics"),
						dsl.I(
							dsl.Class("fas fa-chart-line"),
						),
						dsl.Span(
							dsl.Text("Analytics"),
						),
					),
					dsl.A(
						dsl.Class("nav-item"),
						dsl.AttrHXGET("/users"),
						dsl.I(
							dsl.Class("fas fa-users"),
						),
						dsl.Span(
							dsl.Text("Users"),
						),
					),
				),

				dsl.Div(
					dsl.Class("sidebar-bottom"),
					dsl.Button(
						dsl.Class("theme-toggle-btn"),
						dsl.AttrID("desktopThemeToggle"),
						dsl.I(
							dsl.Class("fas fa-moon"),
						),
						dsl.Span(
							dsl.Text("Dark Mode"),
						),
						dsl.Div(
							dsl.Class("toggle-track"),
						),
						dsl.Div(
							dsl.Class("toggle-thumb"),
						),
					),
				),

				dsl.Div(
					dsl.Class("user-info-sidebar"),
					dsl.Div(
						dsl.Class("user-avatar"),
						dsl.AttrID("userAvatar"),
						dsl.Text("JD"),
					),
					dsl.Div(
						dsl.Div(
							dsl.AttrWithValue("style", "font-weight:600;"),
							dsl.AttrID("userNameDisplay"),
							dsl.Text("John Doe"),
						),
						dsl.Div(
							dsl.AttrWithValue("style", "font-size:0.75rem; color:var(--sidebar-text-muted);"),
							dsl.Text("Admin"),
						),
					),
				),

				dsl.Button(
					dsl.Class("logout-btn"),
					dsl.AttrHXPOST("/logout"),
					dsl.I(
						dsl.Class("fas fa-sign-out-alt"),
					),
					dsl.Text("Logout"),
				),
			),

			dsl.Div(
				dsl.Class("sidebar-overlay"),
				dsl.AttrID("sidebarOverlay"),
			),

			dsl.Main(
				dsl.Class("main-content"),
				dsl.AttrID("main-content"),

				content,
			),

			_MobileBottomNavigation,
		),
	)
}
