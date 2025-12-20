package main

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/handler"
	"DewaSRY/go-microservices/services/api-gateway/internal/service"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/lib"
	"DewaSRY/go-microservices/shared/logger"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	serviceName = "API_GATEWAY"
	PORT        = env.GetString("PORT", "8081")
	rabbitMqURI = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	rabbitmq, err := messaging.NewRabbitMQManager(ctx, rabbitMqURI)
	logger := logger.New()

	defer func() {
		cancel()
		rabbitmq.Close()
	}()

	if err != nil {
		logger.Error("failed_to_connect_to_rabbitmq", err, map[string]interface{}{
			"service_name": serviceName,
		})
		return
	}

	// INIT
	mux := http.NewServeMux()
	con_manager := lib.NewConnectionManager()
	httpHandler := handler.NewHttpHandler()
	tripService := service.NewRideShareService(rabbitmq)
	wsHandler := handler.NewWsHandler(con_manager, rabbitmq, tripService, logger)

	//REGISTER HANDLER
	mux.HandleFunc("GET /health", httpHandler.GetHealthCheck)
	mux.HandleFunc("POST /trip/preview", httpHandler.PostTripPreview)
	mux.HandleFunc("POST /trip/start", httpHandler.PostStartTrip)
	mux.HandleFunc("/ws/connect", wsHandler.WsHandleStartConnection)

	// wrap the handler
	warpHandler := middleware.WithCORS(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: warpHandler,
	}

	go func() {
		logger.Info(fmt.Sprintf("starting_app:%s", serviceName), map[string]interface{}{
			"service_name": serviceName,
		})
		if err := server.ListenAndServe(); err != nil {
			logger.Error(fmt.Sprintf("failed_to_run:%s", serviceName), err, map[string]interface{}{
				"service_name": serviceName,
			})
		}
	}()

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGALRM)
	<-quite
	logger.Info("shout_down_the_server", map[string]interface{}{
		"message": "shout_down_the_server",
	})

	// shout down the server gracefully
	func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("failed_to_shout_down", err, map[string]interface{}{
				"message": "failed_to_shout_down",
			})
			server.Close()
		}

		logger.Info("gracefully_Shout_down", map[string]interface{}{
			"message": "gracefully_Shout_down",
		})
	}()

}
