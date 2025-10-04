package main

import (
	"DewaSRY/go-microservices/services/trip-service/internal/handlers"
	"DewaSRY/go-microservices/services/trip-service/internal/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serverName = "trip_service"
	PORT       = env.GetString("PORT", "8083")
)

func main() {
	log.Printf("start_service:%s\n", serverName)

	tripRepo := repository.NewInMemoryTripRepository()
	tripService := service.NewTripService(tripRepo)
	tripHttpHandler := handlers.NewHttpHandler(tripService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /preview", tripHttpHandler.GetTripPreview)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	go func() {
		log.Printf("success_run_service:%s\n", serverName)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen : %s\n", err)
		}
	}()

	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
		<-quit
		log.Println("shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("server_forced_to_shutdown:%v", err)
		}
	}

	log.Println("server_exiting_gracefully")
}
