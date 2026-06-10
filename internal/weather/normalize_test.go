package weather_test

import (
	"testing"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/stretchr/testify/assert"
)

func TestCityNormalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"capitalize", "Moscow", "moscow"},
		{"trim spaces", "  Moscow ", "moscow"},
		{"lowercase", "MOSCOW", "moscow"},
		{"multiple spaces", "St  Petersburg", "st petersburg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := weather.CityNormalize(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
