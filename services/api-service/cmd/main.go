package main

import (
	"DewaSRY/go-microservices/services/api-service/internal/events"
	"DewaSRY/go-microservices/services/api-service/internal/handler"
	"DewaSRY/go-microservices/services/api-service/internal/repository"
	"DewaSRY/go-microservices/services/api-service/internal/service"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/logger"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	// GrpcAddr     = env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	// appURL       = env.GetString("APP_URL", "http://localhost:3000")
	postgres_uri = env.GetString("POSTGRES_URI", "postgres://postgres:postgres@postgres:5432/riderdb?sslmode=disable")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	logger := logger.New()

	rabbitmq, err := messaging.NewRabbitMQManager(ctx, rabbitMqURI)
	if err != nil {
		logger.Error("failed_to_connect_to_rabbitmq", err, map[string]interface{}{
			"service_name": "API_SERVICE",
		})
		return
	}

	db, err := db.NewPostgresManager(ctx, postgres_uri)
	if err != nil {
		logger.Error("failed_to_connect_to_db", err, map[string]interface{}{
			"event": "failed_to_connect_to_db",
			"error": err,
		})

		cancel()
		return
	}

	userRepository := repository.NewUserRepository(db)
	tripFlowRepository := repository.NewTripFlowRepository(db)

	userService := service.NewUserService(rabbitmq, userRepository, tripFlowRepository)
	tripFlowService := service.NewTripFlowService(rabbitmq, tripFlowRepository)
	osrmService := service.NewOsrmIntegrationService()

	userEventHandler := handler.NewUserEventHandler(userService, rabbitmq, logger)
	tripFlowEventHandler := handler.NewTripFlowEventHandler(rabbitmq, tripFlowService, userService, osrmService, logger)

	userEventConsumer := events.NewUserConsumer(rabbitmq, userEventHandler)
	tripFlowConsumer := events.NewTripFlowConsumer(rabbitmq, tripFlowEventHandler)

	go func() {
		if err := userEventConsumer.Listen(); err != nil {
			logger.Error("failed_to_listen_user_event_flow_queue", err, map[string]interface{}{
				"event": "failed_to_listen_user_event_flow_queue",
				"error": err,
			})
			cancel()
		}
	}()

	go func() {
		if err := tripFlowConsumer.Listen(); err != nil {
			logger.Error("failed_to_listen_trip_event_flow_queue", err, map[string]interface{}{
				"event": "failed_to_listen_trip_event_flow_queue",
				"error": err,
			})
			cancel()
		}
	}()

	logger.Info("server is starting", map[string]interface{}{
		"event": "api_service_started",
		"msg":   "API Service is up and running",
	})

	defer func() {
		cancel()
		rabbitmq.Close()
	}()

	<-ctx.Done()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		<-sigCh
		cancel()
	}()
	logger.Info("server_is_gracefully_down", map[string]interface{}{
		"event": "server_is_gracefully_down",
		"msg":   "Api Service is down",
	})

}
