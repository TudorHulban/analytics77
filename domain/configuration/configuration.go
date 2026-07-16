package configuration

// Configuration holds the main application configuration mapping.
type Configuration struct {
	Debug  ConfigDebug  `json:"debug"`
	Server ConfigServer `json:"server"`
}

// ConfigServer stores the configuration for the primary application server.
type ConfigServer struct {
	Port int `json:"port"`
}

// ConfigDebug stores the configuration for the debugging interface.
type ConfigDebug struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}
