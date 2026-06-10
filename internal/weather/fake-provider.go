package weather

import "context"

type FakeProvider struct {
	Snapshot WeatherSnapshot
	Err      error
}

func (f *FakeProvider) Current(ctx context.Context, city string) (WeatherSnapshot, error) {
	return f.Snapshot, f.Err
}
