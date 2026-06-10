package main

import (
	"context"
	_ "fmt"
	"log"
	_ "os/signal"
	_ "syscall"
	"time"
	_ "time"

	inmemorycache "github.com/ScaryOrange/go-weather-snapshot-service/internal/cache/in-memory-cache"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/config"
	myHttp "github.com/ScaryOrange/go-weather-snapshot-service/internal/http"
	myMetrics "github.com/ScaryOrange/go-weather-snapshot-service/internal/metrics"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/storage/postgres"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	provider := weather.NewOpenMeteoClient(cfg.GeocodingEndpoint, cfg.ForecastEndpoint)

	dbPool, err := pgxpool.New(ctx, cfg.DB)
	if err != nil {
		log.Fatalf("Database connetcion error: %v", err)
	}
	defer dbPool.Close()
	storage := postgres.NewStorage(dbPool)

	cache := inmemorycache.NewMemoryCache(time.Duration(cfg.TTL) * time.Minute)

	weatherMetrics := myMetrics.NewWeatherMetrics()

	handler := myHttp.NewHandler(provider, storage, cache, weatherMetrics)

	myHttp.RunServer(cfg.Port, handler)

}
