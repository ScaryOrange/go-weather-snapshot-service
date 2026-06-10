package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port              string `env:"PORT" env-required:"true"`
	Host              string `env:"HOST"`
	DB                string `env:"DB"`
	GeocodingEndpoint string `env:"GEOCODING_ENDPOINT"`
	ForecastEndpoint  string `env:"FORECAST_ENDPOINT"`
	TTL               int64  `env:"TTL"`
}

func Load() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
