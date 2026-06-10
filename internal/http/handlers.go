package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/cache"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/metrics"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/storage"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
)

type history struct {
	Items []weather.WeatherSnapshot `json:"items"`
}

type handler struct {
	provider weather.WeatherProvider
	storage  storage.WeatherStorage
	cache    cache.WeatherCache
	metrics  *metrics.WeatherMetrics
}

func NewHandler(
	provider weather.WeatherProvider,
	storage storage.WeatherStorage,
	cache cache.WeatherCache,
	metrics *metrics.WeatherMetrics,
) *handler {
	return &handler{provider: provider, storage: storage, cache: cache, metrics: metrics}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondCurrentWeather(w http.ResponseWriter, snapshot weather.WeatherSnapshot) {
	log.Printf("STATUS %d: http.GetCurrentWeather", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(snapshot)
}

func (h *handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (h *handler) GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if len(city) < 2 {
		if h.metrics != nil {
			h.metrics.ErrorInc()
		}
		log.Printf("ERROR %d: city=%s, city lenght must be more than 2", http.StatusBadRequest, city)
		respondWithError(w, http.StatusBadRequest, "city is required and must be at least 2 characters")
		return
	}
	city = weather.CityNormalize(city)
	if snapshot, ok := h.cache.Get(city); ok == true {
		if h.metrics != nil {
			h.metrics.CacheHitInc()
		}
		respondCurrentWeather(w, snapshot)
		return
	}
	if h.metrics != nil {
		h.metrics.CacheMissInc()
	}

	snapshot, err := h.provider.Current(r.Context(), city)
	if err != nil {
		if h.metrics != nil {
			h.metrics.ErrorInc()
		}
		log.Printf("ERROR %d: fiald to provider.Current, %v", http.StatusInternalServerError, err)
		respondWithError(w, http.StatusInternalServerError, "fiald to fetch weather")
		return
	}

	if err := h.storage.Save(r.Context(), snapshot); err != nil {
		log.Printf("ERROR %d: faild to storage.Save, %v", http.StatusInternalServerError, err)
	}

	h.cache.Set(city, snapshot)
	respondCurrentWeather(w, snapshot)
}

func (h *handler) GetHistoryWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if len(city) < 2 {
		log.Printf("ERROR %d: city=%s, city lenght must be more than 2", http.StatusBadRequest, city)
		respondWithError(w, http.StatusBadRequest, "city is required and must be at least 2 characters")
		return
	}
	city = weather.CityNormalize(city)

	var limit int
	tmpLimit := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(tmpLimit)
	if err != nil {
		log.Printf("NOTICE: invalid limit \"%s\"", tmpLimit)
		limit = 10
	}

	items, err := h.storage.History(r.Context(), city, limit)
	if err != nil {
		log.Printf("ERROR: databse error %v", err)
	}
	if len(items) == 0 {
		items = []weather.WeatherSnapshot{}
	}
	response := history{Items: items}

	log.Printf("STATUS %d: http.GetHistoryWeather", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
