package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type WeatherMetrics struct {
	CurrentRequests *prometheus.CounterVec
	CacheHits       prometheus.Counter
	CacheMisses     prometheus.Counter
	registry        *prometheus.Registry
}

func NewWeatherMetrics() *WeatherMetrics {
	registry := prometheus.NewRegistry()
	m := &WeatherMetrics{
		CurrentRequests: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "total_current_requests_count",
			},
			[]string{"result"},
		),
		CacheHits: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_hits_count",
			},
		),
		CacheMisses: promauto.NewCounter(
			prometheus.CounterOpts{
				Name: "cache_misses_count",
			},
		),
		registry: registry,
	}
	registry.MustRegister(m.CurrentRequests, m.CacheHits, m.CacheMisses)

	return m
}

func (m *WeatherMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

func (m *WeatherMetrics) CacheHitInc() {
	m.CacheHits.Inc()
	m.CurrentRequests.WithLabelValues("cache_hit").Inc()
}

func (m *WeatherMetrics) CacheMissInc() {
	m.CacheMisses.Inc()
	m.CurrentRequests.WithLabelValues("cache_miss").Inc()
}

func (m *WeatherMetrics) ErrorInc() {
	m.CurrentRequests.WithLabelValues("error").Inc()
}
