package configuration

type IConfiguration interface {
	Get() map[string]any
}
