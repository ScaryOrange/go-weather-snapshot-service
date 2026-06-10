package inmemorycache_test

import (
	"testing"
	"time"

	inmemorycache "github.com/ScaryOrange/go-weather-snapshot-service/internal/cache/in-memory-cache"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCacheSetAndGet(t *testing.T) {
	c := inmemorycache.NewMemoryCache(5 * time.Minute)

	snapshot := weather.WeatherSnapshot{City: "Moscow", TemperatureCelsius: 15.5}
	c.Set("moscow", snapshot)

	cached, ok := c.Get("moscow")
	assert.True(t, ok)
	assert.Equal(t, snapshot.City, cached.City)
	assert.Equal(t, snapshot.TemperatureCelsius, cached.TemperatureCelsius)
}

func TestInMemoryCacheGetExpired(t *testing.T) {
	cache := inmemorycache.NewMemoryCache(1 * time.Second)
	cache.Set("moscow", weather.WeatherSnapshot{})
	time.Sleep(2 * time.Second)
	_, ok := cache.Get("moscow")
	assert.False(t, ok)
}

func TestMemorryCacheGetMiss(t *testing.T) {
	cache := inmemorycache.NewMemoryCache(1 * time.Second)
	_, ok := cache.Get("Moscow")
	assert.False(t, ok)
}

func TestMemoryCacheConcurencyAccess(t *testing.T) {
	c := inmemorycache.NewMemoryCache(5 * time.Minute)
	done := make(chan bool)

	go func() {
		for range 1000 {
			c.Set("key", weather.WeatherSnapshot{City: "Moscow"})
		}
		done <- true
	}()

	go func() {
		for range 1000 {
			c.Get("key")
		}
		done <- true
	}()

	<-done
	<-done
}
