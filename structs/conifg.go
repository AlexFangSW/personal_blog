package structs

type ServerSetting struct {
	Port   string `json:"port"`
	Prefix string `json:"prefix"`
}

type LoggerSetting struct {
	Level string `json:"level"`
}

type Config struct {
	Server ServerSetting `json:"server"`
	Logger LoggerSetting `json:"logger"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerSetting{
			Port:   ":8080",
			Prefix: "/api/v1",
		},
		Logger: LoggerSetting{
			Level: "INFO",
		},
	}
}
