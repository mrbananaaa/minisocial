package config

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Database DatabaseConfig
	NATS     NATSConfig
	Logger   LoggerConfig
}

type AppConfig struct {
	Env string
}

type HTTPConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type NATSConfig struct {
	URL string
}

type LoggerConfig struct {
	Level string
}

func (c *Config) IsDev() bool {
	return c.App.Env == "development"
}
