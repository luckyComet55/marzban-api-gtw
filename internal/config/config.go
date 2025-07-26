package config

type AppConfig struct {
	Username       string `env:"ADMIN_USERNAME, required"`
	Password       string `env:"ADMIN_PASSWORD, required"`
	MarzbanBaseUrl string `env:"BASE_URL, required"`
	Port           int    `env:"PORT, default=8343"`
}
