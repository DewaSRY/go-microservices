package main

import (
	"DewaSRY/go-microservices/services/api-service/internal/events"
	"DewaSRY/go-microservices/services/api-service/internal/handler"
	"DewaSRY/go-microservices/services/api-service/internal/repository"
	"DewaSRY/go-microservices/services/api-service/internal/service"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	GrpcAddr     = env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI  = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	appURL       = env.GetString("APP_URL", "http://localhost:3000")
	postgres_uri = env.GetString("POSTGRES_URI", "postgres://postgres:postgres@postgres:5432/riderdb?sslmode=disable")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		<-sigCh
		cancel()
	}()

	rabbitmq, err := messaging.NewRabbitMQManager(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	db, err := db.NewPostgresManager(ctx, postgres_uri)
	if err != nil {
		log.Printf("failed_to_make_db_connection:%v", err)
		cancel()
		return
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(rabbitmq, userRepository)
	userEventHandler := handler.NewUserEventHandler(userService)

	userEventConsumer := events.NewUserConsumer(*rabbitmq, userEventHandler)

	go func() {
		if err := userEventConsumer.Listen(); err != nil {
			log.Printf("failed_to_listen_service:%v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting_down_payment_service")

}
