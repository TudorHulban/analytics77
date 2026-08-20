package layouts

import "github.com/TudorHulban/hxgo/dsl"

var _LinkCSSStyles = dsl.Link(
	dsl.Rel("stylesheet"),
	dsl.Href("/public/styles.css"),
)

var _JSCore = dsl.Script(
	dsl.Src("/public/hxgo_core_ws.js"),
)

var _JSCache = dsl.Script(
	dsl.Src("/public/hxgo_plugin_cache.js"),
)

var _JSListeners = dsl.Script(
	dsl.Src("/public/hxgo_plugin_listeners.js"),
)

var _JSUI = dsl.Script(
	dsl.Src("/public/hxgo_plugin_ui.js"),
)

var _JSValidation = dsl.Script(
	dsl.Src("/public/hxgo_plugin_validation.js"),
)

var _JSChart = dsl.Script(
	dsl.Src("https://cdn.jsdelivr.net/npm/chart.js@4.4.1/dist/chart.umd.min.js"),
)

var _JS = []dsl.Node{
	_JSCore,
	_JSCache,
	_JSListeners,
	_JSUI,
	_JSValidation,
	_JSChart,
}
