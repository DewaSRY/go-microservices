package main

import (
	"DewaSRY/go-microservices/services/driver-service/internal/events"
	"DewaSRY/go-microservices/services/driver-service/internal/handler"
	"DewaSRY/go-microservices/services/driver-service/internal/repository"
	"DewaSRY/go-microservices/services/driver-service/internal/service"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

var (
	GrpcAddr      = ":9092"
	amqpUrlString = env.GetString("AMQP_URL", "amqp://guess:guess@rabbitmq:5672/")
	postgres_uri  = env.GetString("POSTGRES_URI", "postgres://postgres:postgres@postgres:5432/riderdb?sslmode=disable")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	lis, err := net.Listen("tcp", GrpcAddr)
	if err != nil {
		log.Fatalf("failed_to_listen: %v", err)
		return
	}

	message_manager, err := messaging.NewRabbitMQManager(ctx, amqpUrlString)
	if err != nil {
		log.Fatalf("failed_to_listen: %v", err)
		return
	}
	defer message_manager.Close()

	db, err := db.NewPostgresManager(ctx, postgres_uri)
	if err != nil {
		log.Printf("failed_to_make_db_connection:%v", err)
		cancel()
		return
	}

	driverRepo := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(driverRepo)
	grpcServer := grpc.NewServer()
	handler.NewGrpcHandler(grpcServer, driverService)
	event_handler := handler.NewTripEventHandler(driverService, message_manager)
	tripListener := events.NewTripConsumer(message_manager, event_handler)

	go func() {
		if err := tripListener.Listen(); err != nil {
			log.Fatalf("failed_to_make_listener_to_trip_event:%v", err)
		}
	}()

	log.Printf("starting_grpc_server_driverService_on_port:%s", lis.Addr().String())

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed_to_serve:%s", err)
			cancel()
		}
	}()

	defer func() {
		message_manager.Close()
		cancel()
	}()

	go func() {
		sigCn := make(chan os.Signal, 1)
		signal.Notify(sigCn, os.Interrupt, syscall.SIGTERM)
		<-sigCn
		cancel()
	}()
	<-ctx.Done()
	log.Println("shutting_down_the_server..")
	grpcServer.GracefulStop()
}
