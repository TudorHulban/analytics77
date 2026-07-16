package appanalytics

func (a *App) Stop() {
	a.fnFreeResources()
}
