package http

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RunServer(port string, myHandler *handler) {
	router := chi.NewRouter()
	router.Get("/api/v1/health", myHandler.Health)
	router.Get("/api/v1/weather/current", myHandler.GetCurrentWeather)
	router.Get("/api/v1/weather/history", myHandler.GetHistoryWeather)
	router.Get("/metrics", myHandler.metrics.Handler().ServeHTTP)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
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
