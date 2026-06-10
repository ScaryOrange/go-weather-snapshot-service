package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	inmemorycache "github.com/ScaryOrange/go-weather-snapshot-service/internal/cache/in-memory-cache"
	myHttp "github.com/ScaryOrange/go-weather-snapshot-service/internal/http"
	fakestorage "github.com/ScaryOrange/go-weather-snapshot-service/internal/storage/fake-storage"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/stretchr/testify/assert"
)

func TestGetCurrentWeather_Success(t *testing.T) {
	provider := &weather.FakeProvider{
		Snapshot: weather.WeatherSnapshot{
			City:               "Moscow",
			Provider:           "open-meteo",
			TemperatureCelsius: 15.2,
			WindSpeed:          3.5,
			ObservedAt:         time.Now(),
		},
	}
	cache := inmemorycache.NewMemoryCache(5 * time.Second)
	storage := &fakestorage.FakeStorage{}
	handler := myHttp.NewHandler(provider, storage, cache, nil)

	req := httptest.NewRequest("GET", "/api/v1/weather/current?city=Moscow", nil)
	rr := httptest.NewRecorder()

	handler.GetCurrentWeather(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp weather.WeatherSnapshot
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Cached)
	assert.Equal(t, "Moscow", resp.City)
}

func TestGetCurrentWeather_CacheHit(t *testing.T) {
	provider := &weather.FakeProvider{
		Snapshot: weather.WeatherSnapshot{City: "Moscow", TemperatureCelsius: 22.5},
	}
	storage := &fakestorage.FakeStorage{}
	cache := inmemorycache.NewMemoryCache(5 * time.Minute)
	handler := myHttp.NewHandler(provider, storage, cache, nil)

	req1 := httptest.NewRequest("GET", "/api/v1/weather/current?city=Moscow", nil)
	rr1 := httptest.NewRecorder()
	handler.GetCurrentWeather(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	req2 := httptest.NewRequest("GET", "/api/v1/weather/current?city=Moscow", nil)
	rr2 := httptest.NewRecorder()
	handler.GetCurrentWeather(rr2, req2)

	var resp weather.WeatherSnapshot
	json.Unmarshal(rr2.Body.Bytes(), &resp)
	assert.True(t, resp.Cached)
}

func TestGetCurrentWeather_ProviderError(t *testing.T) {
	provider := &weather.FakeProvider{Err: errors.New("external API error")}
	storage := &fakestorage.FakeStorage{}
	cache := inmemorycache.NewMemoryCache(5 * time.Minute)
	handler := myHttp.NewHandler(provider, storage, cache, nil)

	req := httptest.NewRequest("GET", "/api/v1/weather/current?city=Moscow", nil)
	rr := httptest.NewRecorder()
	handler.GetCurrentWeather(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestGetCurrentWeather_InvalidCity(t *testing.T) {
	provider := &weather.FakeProvider{}
	storage := &fakestorage.FakeStorage{}
	cache := inmemorycache.NewMemoryCache(5 * time.Minute)
	handler := myHttp.NewHandler(provider, storage, cache, nil)

	req := httptest.NewRequest("GET", "/api/v1/weather/current?city=", nil)
	rr := httptest.NewRecorder()
	handler.GetCurrentWeather(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetHistory_Success(t *testing.T) {
	now := time.Now()
	expectedHistory := []weather.WeatherSnapshot{
		{City: "Moscow", TemperatureCelsius: 20.5, ObservedAt: now},
		{City: "Moscow", TemperatureCelsius: 21.0, ObservedAt: now.Add(-time.Hour)},
	}

	storage := &fakestorage.FakeStorage{Snapshots: expectedHistory}
	handler := myHttp.NewHandler(nil, storage, nil, nil)

	req := httptest.NewRequest("GET", "/api/v1/weather/history?city=Moscow&limit=10", nil)
	rr := httptest.NewRecorder()

	handler.GetHistoryWeather(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Items []weather.WeatherSnapshot `json:"items"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Items, 2)
	assert.Equal(t, "Moscow", response.Items[0].City)
}

func TestGetHistory_Empty(t *testing.T) {
	storage := &fakestorage.FakeStorage{Snapshots: []weather.WeatherSnapshot{}}
	handler := myHttp.NewHandler(nil, storage, nil, nil)

	req := httptest.NewRequest("GET", "/api/v1/weather/history?city=Moscow&limit=10", nil)
	rr := httptest.NewRecorder()

	handler.GetHistoryWeather(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response struct {
		Items []weather.WeatherSnapshot `json:"items"`
	}
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Empty(t, response.Items)
	assert.Contains(t, rr.Body.String(), `"items":[]`)
}
