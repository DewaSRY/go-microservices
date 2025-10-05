package main

import (
	"DewaSRY/go-microservices/services/trip-service/internal/handlers"
	"DewaSRY/go-microservices/services/trip-service/internal/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

// TODO clean up
var (
	serverName = "trip_service"
	PORT       = env.GetString("PORT", "9093")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tripRepo := repository.NewInMemoryTripRepository()
	tripService := service.NewTripService(tripRepo)
	tripFareService := service.NewTripFareService(tripRepo)
	// tripHandler := handlers.NewHttpHandler(tripService)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("failed_to_listen:%v", err)
	}

	grpcService := grpcserver.NewServer()
	handlers.NewGRPCHandler(grpcService, tripService, tripFareService)

	go func() {
		log.Printf("success_run_service:%s on port %s\n", serverName, fmt.Sprintf(":%s", PORT))
		if err := grpcService.Serve(lis); err != nil {
			log.Fatalf("Listen : %s\n", err)
		}
	}()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	<-ctx.Done()
	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGALRM)
		<-quit
		log.Println("shutting down server...")
		grpcService.Stop()
		log.Println("server_exiting_gracefully")
	}
}
