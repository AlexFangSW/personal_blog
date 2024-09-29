package config

type ServerSetting struct {
	Port            int    `json:"port"`
	Prefix          string `json:"prefix"`
	ShutdownTimeout int    `json:"shutdownTimeout"`
}

type LoggerSetting struct {
	Level string `json:"level"`
}

type DBSetting struct {
	DSNURL      string `json:"dsnURL"`
	Timeout     int    `json:"timeout"`
	Connections int    `json:"connections"`
}

type JWTSetting struct {
	Issuer string `json:"issuer"`
	// hour
	Expire int    `json:"expire"`
	Secret string `json:"secret"`
}

type LoginSetting struct {
	RateLimit int `json:"rateLimit"` // request per second
}

type Config struct {
	Server ServerSetting `json:"server"`
	Logger LoggerSetting `json:"logger"`
	DB     DBSetting     `json:"db"`
	JWT    JWTSetting    `json:"jwt"`
	Login  LoginSetting  `json:"login"`
}

func NewConfig() *Config {
	return &Config{
		Server: ServerSetting{
			Port:            8080,
			Prefix:          "/api/v1",
			ShutdownTimeout: 30,
		},
		Logger: LoggerSetting{
			Level: "INFO",
		},
		DB: DBSetting{
			DSNURL:      "./example.db",
			Timeout:     30,
			Connections: 10,
		},
		JWT: JWTSetting{
			Issuer: "alexfangsw",
			Expire: 6,
		},
		Login: LoginSetting{
			RateLimit: 1,
		},
	}
}
