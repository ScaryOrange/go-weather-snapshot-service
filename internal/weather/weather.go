package weather

import (
	"context"
	"strings"
	"time"
)

type WeatherProvider interface {
	Current(ctx context.Context, city string) (WeatherSnapshot, error)
}

type WeatherSnapshot struct {
	City               string    `json:"city"`
	Provider           string    `json:"provider"`
	TemperatureCelsius float64   `json:"temperature_celsius"`
	WindSpeed          float64   `json:"wind_speed"`
	ObservedAt         time.Time `json:"observed_at"`
	RawPayload         []byte    `json:"-"`
	Cached             bool      `json:"cached"`
}

func CityNormalize(city string) string {
	normCity := strings.Trim(city, " ")
	normCity = strings.Trim(normCity, "\n")
	normCity = strings.ToLower(normCity)
	return strings.Join(strings.Fields(normCity), " ")
}
