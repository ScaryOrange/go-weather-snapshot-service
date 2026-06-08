package storage

import (
	"context"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
)

type WeatherStorage interface {
	Save(ctx context.Context, snapshot weather.WeatherSnapshot) error
	History(ctx context.Context, city string, limit int) ([]weather.WeatherSnapshot, error)
}
