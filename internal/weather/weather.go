package weather

import (
	"context"
	"strings"
)

type WeatherProvider interface {
	Current(ctx context.Context, city string) (WeatherSnapshot, error)
}

type WeatherSnapshot struct {
	City               string
	Provider           string
	TemperatureCelsius float64
	WindSpeed          float64
	ObservedAt         string
	RawPayload         []byte
}

func CityNormalize(city string) string {
	normCity := strings.Trim(city, " ")
	normCity = strings.Trim(normCity, "\n")
	normCity = strings.ToLower(normCity)
	return normCity
}
