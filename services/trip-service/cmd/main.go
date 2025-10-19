package main

import (
	"DewaSRY/go-microservices/services/trip-service/internal/events"
	"DewaSRY/go-microservices/services/trip-service/internal/handlers"
	"DewaSRY/go-microservices/services/trip-service/internal/repository"
	"DewaSRY/go-microservices/services/trip-service/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/tracing"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "google.golang.org/grpc"
)

var (
	serverName      = "trip_service"
	PORT            = env.GetString("PORT", "9093")
	amqp_url_string = env.GetString("AMQP_URL", "amqp://guess:guess@rabbitmq:5672/")
)

func main() {
	tracerCfg := tracing.Config{
		ServiceName:    "trip-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatalf("failed_to_listen:%v", err)
	}

	conn, err := messaging.NewRabbitMQManager(amqp_url_string)
	log.Print("start connect to rabbitmq ")
	if err != nil {
		log.Printf("failed_to_make_connection:%v", err)
		cancel()
		return
	}

	defer conn.Close()

	grpcService := grpcserver.NewServer()

	tripRepo := repository.NewInMemoryTripRepository()
	tripService := service.NewTripService(tripRepo)
	tripFareService := service.NewTripFareService(tripRepo)
	tripDriverEventHandler := handlers.NewTripDriverEventHandler(conn, tripService, tripFareService)
	tripPaymentEventHandler := handlers.NewTripPaymentEventHandler(conn, tripService)

	tripEventPublisher := events.NewTripEventPublisher(conn)
	tripDriverEventConsumer := events.NewDriverConsumer(conn, tripDriverEventHandler)
	tripPaymentEventConsumer := events.NewPaymentConsumer(conn, tripPaymentEventHandler)

	// tripHandler := handlers.NewHttpHandler(tripService)
	handlers.NewGRPCHandler(grpcService, tripService, tripFareService, tripEventPublisher)

	go func() {
		log.Printf("success_run_service:%s on port %s\n", serverName, fmt.Sprintf(":%s", PORT))
		if err := grpcService.Serve(lis); err != nil {
			log.Fatalf("Listen : %s\n", err)
		}
	}()

	go func() {
		if err := tripDriverEventConsumer.Listen(); err != nil {
			log.Fatalf("failed_to_make_listener_to_driven_event:%v", err)
		}
	}()

	go func() {
		if err := tripPaymentEventConsumer.Listen(); err != nil {
			log.Fatalf("failed_to_make_listener_to_payment_event:%v", err)
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
