package sconfiguration

type ServiceConfiguration struct {
	pathConfiguration string
}

func NewServiceConfiguration(pathConfiguration string) (*ServiceConfiguration, error) {
	return &ServiceConfiguration{
			pathConfiguration: pathConfiguration,
		},
		nil
}

func (*ServiceConfiguration) Get() map[string]any {
	return map[string]any{}
}
