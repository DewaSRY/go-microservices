package main

import (
	"DewaSRY/go-microservices/services/payment-service/internal/events"
	"DewaSRY/go-microservices/services/payment-service/internal/handler"
	"DewaSRY/go-microservices/services/payment-service/internal/service"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"DewaSRY/go-microservices/shared/env"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/tracing"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	GrpcAddr    = env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	appURL      = env.GetString("APP_URL", "http://localhost:3000")
)

func main() {
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
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		<-sigCh
		cancel()
	}()

	// Stripe config
	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: env.GetString("STRIPE_SECRET_KEY", ""),
		SuccessURL:      env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success"),
		CancelURL:       env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel"),
	}

	if stripeCfg.StripeSecretKey == "" {
		log.Fatalf("STRIPE_SECRET_KEY is not set")
		return
	}

	rabbitmq, err := messaging.NewRabbitMQManager(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()
	log.Println("Starting RabbitMQ connection")

	paymentProcessor := service.NewPaymentProcessService(stripeCfg)
	service := service.NewPaymentService(paymentProcessor)
	tripHandler := handler.NewTripEventHandler(*rabbitmq, service)

	tripConsumer := events.NewTripConsumer(*rabbitmq, tripHandler)

	go func() {
		if err := tripConsumer.Listen(); err != nil {
			log.Printf("failed_to_listen_service:%v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	log.Println("shutting_down_payment_service")

}
