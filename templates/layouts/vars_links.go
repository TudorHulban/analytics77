package layouts

import "github.com/TudorHulban/hxgo/dsl"

var _LinkCSSStyles = dsl.Link(
	dsl.Rel("stylesheet"),
	dsl.Href("/public/styles.css"),
)

var _LinkCSSFontAwesome = dsl.Link(
	dsl.Rel("stylesheet"),
	dsl.Href("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0-beta3/css/all.min.css"),
)

var _JSDashboard = dsl.Script(
	dsl.Src("/public/dashboard.js"),
)

var _JSCore = dsl.Script(
	dsl.Src("/public/hxgo/js/hxgo_core_ws.js"),
)

var _JSCache = dsl.Script(
	dsl.Src("/public/hxgo/js/hxgo_plugin_cache.js"),
)

var _JSListeners = dsl.Script(
	dsl.Src("/public/hxgo/js/hxgo_plugin_listeners.js"),
)

var _JSUI = dsl.Script(
	dsl.Src("/public/hxgo/js/hxgo_plugin_ui.js"),
)

var _JSValidation = dsl.Script(
	dsl.Src("/public/hxgo/js/hxgo_plugin_validation.js"),
)

var _JSChart = dsl.Script(
	dsl.Src("https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"),
)

var _JS = []dsl.Node{
	_JSDashboard,

	_JSCore,
	_JSCache,
	_JSListeners,
	_JSUI,
	_JSValidation,
	_JSChart,
}
