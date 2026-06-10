package cache

import (
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
)

type WeatherCache interface {
	Get(city string) (weather.WeatherSnapshot, bool)
	Set(city string, snapshot weather.WeatherSnapshot) // в тз передается ttl, но принято решение реализовать в структуре кеша поле .ttl
}
