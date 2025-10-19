package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"DewaSRY/go-microservices/services/api-gateway/internal/handler"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/lib"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/middleware"
	"DewaSRY/go-microservices/shared/tracing"
)

var (
	serviceName = "API_GATEWAY"
	PORT        = env.GetString("PORT", "8081")
	rabbitMqURI = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {

	tracerCfg := tracing.Config{
		ServiceName:    "api-gateway",
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

	rabbitmq, err := messaging.NewRabbitMQManager(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	// INIT
	mux := http.NewServeMux()
	con_manager := lib.NewConnectionManager()
	httpHandler := handler.NewHttpHandler()
	wsHandler := handler.NewWsHandler(con_manager, rabbitmq)

	//REGISTER HANDLER
	mux.Handle("GET /health", tracing.WrapHandlerFunc(httpHandler.GetHealthCheck, "/health"))
	mux.Handle("POST /trip/preview", tracing.WrapHandlerFunc(httpHandler.PostTripPreview, "/trip/preview"))
	mux.Handle("POST /trip/start", tracing.WrapHandlerFunc(httpHandler.PostStartTrip, "/trip/start"))

	mux.HandleFunc("/ws/riders", wsHandler.WsHandleRider)
	mux.HandleFunc("/ws/drivers", wsHandler.WsHandleDriver)

	// wrap the handler
	warpHandler := middleware.WithCORS(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: warpHandler,
	}

	go func() {
		log.Printf("starting_app:%s\n", serviceName)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed_to_run:%s", serviceName)
		}
	}()

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGALRM)
	<-quite
	log.Println("shout down the server")
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("failed_to_shout_down")
		server.Close()
	}

	log.Println("gracefully_Shout_down")

}
