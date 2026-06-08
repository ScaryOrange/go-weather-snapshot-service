package main

import (
	_ "context"
	_ "fmt"
	"log"
	"net/http"
	_ "os/signal"
	_ "syscall"
	_ "time"

	"github.com/ScaryOrange/go-weather-snapshot-service/internal/config"
	myHttp "github.com/ScaryOrange/go-weather-snapshot-service/internal/http"
	"github.com/ScaryOrange/go-weather-snapshot-service/internal/weather"
	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
		return
	}

	provider := weather.NewOpenMeteoClient(cfg.GeocodingEndpoint, cfg.ForecastEndpoint)
	_ = provider
	router.Get("/api/v1/health", myHttp.Health)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}

	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer stop()

	// cfg, _ := config.Load()
	// provider := weather.NewOpenMeteoClient(cfg.GeocodingEndpoint, cfg.ForecastEndpoint)
	// fmt.Print(provider)

	// router := chi.NewRouter()
	// router.Get("/api/v1/health", myHttp.Health)

	// server := http.Server{Addr: ":" + cfg.Port, Handler: router}

	// go func() {
	// 	log.Println("HTTP server listening on :" + cfg.Port)
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		log.Fatalf("server error: %v", err)
	// 	}
	// }()

	// <-ctx.Done()
	// log.Println("Shutdown signal")

	// shutdownCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	// defer cancel()
	// if err := server.Shutdown(shutdownCtx); err != nil {
	// 	log.Printf("HTTP server shutdown error: %v", err)
	// }

}
