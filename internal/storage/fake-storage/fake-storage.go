package fakestorage

import (
	"context"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
)

type FakeStorage struct {
	SaveErr    error
	Snapshots  []weather.WeatherSnapshot
	HistoryErr error
}

func (f *FakeStorage) Save(ctx context.Context, snap weather.WeatherSnapshot) error {
	return f.SaveErr
}

func (f *FakeStorage) History(ctx context.Context, city string, limit int) ([]weather.WeatherSnapshot, error) {
	if f.HistoryErr != nil {
		return nil, f.HistoryErr
	}
	return f.Snapshots, nil
}
