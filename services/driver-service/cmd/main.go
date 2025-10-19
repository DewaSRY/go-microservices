package main

import (
	"DewaSRY/go-microservices/services/driver-service/internal/events"
	"DewaSRY/go-microservices/services/driver-service/internal/handler"
	"DewaSRY/go-microservices/services/driver-service/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/tracing"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

var GrpcAddr = ":9092"

func main() {
	amqpUrlString := env.GetString("AMQP_URL", "amqp://guess:guess@rabbitmq:5672/")
	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "payment-service",
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

	go func() {
		sigCn := make(chan os.Signal, 1)
		signal.Notify(sigCn, os.Interrupt, syscall.SIGTERM)
		<-sigCn
		cancel()
	}()

	lis, err := net.Listen("tcp", GrpcAddr)

	if err != nil {
		log.Fatalf("failed_to_listen: %v", err)
		return
	}

	message_manager, err := messaging.NewRabbitMQManager(amqpUrlString)
	if err != nil {
		log.Fatalf("failed_to_listen: %v", err)
		return
	}
	defer message_manager.Close()

	driver_service := service.NewDriverService()
	grpcServer := grpc.NewServer()
	handler.NewGrpcHandler(grpcServer, driver_service)

	event_handler := handler.NewTripEventHandler(driver_service, message_manager)

	tripListener := events.NewTripConsumer(message_manager, event_handler)
	go func() {
		if err := tripListener.Listen(); err != nil {
			log.Fatalf("failed_to_make_listener_to_trip_event:%v", err)
		}
	}()

	log.Printf("starting_grpc_server_driver_service_on_port:%s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%s", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting_down_the_server..")
	grpcServer.GracefulStop()
}
