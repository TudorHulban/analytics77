package appanalytics

func InitializeTransportRoutes(application *App) {
	application.transportHTTP.Get(
		RoutesAll,
		application.HandlerViewRegistry,
	)
}
