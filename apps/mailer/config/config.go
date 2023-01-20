package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Addr          string `json:"addr" env:"MAILER_ADDR" env-required:"true"`
	Password      string `json:"password" env:"MAILER_PASSWORD" env-required:"true"`
	From          string `json:"from" env:"MAILER_FROM" env-required:"true"`
	Port          string `json:"port" env:"MAILER_PORT" env-required:"true"`
	ServerAddress string `json:"server_address" env:"SERVER_ADDRESS" env-default:":8080"`
	LogLevel      string `json:"log_level" env:"LOG_LEVEL" env-default:"debug"`
}

func NewConfig(fromFile bool) (*Config, error) {
	var cfg Config

	if fromFile {
		if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
			return nil, err
		}
		return &cfg, nil
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
