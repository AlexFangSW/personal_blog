package structs

type ServerSetting struct {
	Port   string `json:"port"`
	Prefix string `json:"prefix"`
}

type LoggerSetting struct {
	Level string `json:"level"`
}

type DBSetting struct {
	DSNURL  string `json:"dsnURL"`
	Timeout int    `json:"timeout"`
}

type Config struct {
	Server ServerSetting `json:"server"`
	Logger LoggerSetting `json:"logger"`
	DB     DBSetting     `json:"db"`
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
		DB: DBSetting{
			DSNURL:  "./example.db",
			Timeout: 30,
		},
	}
}
