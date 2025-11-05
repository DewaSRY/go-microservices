package main

import (
	"DewaSRY/go-microservices/services/payment-service/internal/events"
	"DewaSRY/go-microservices/services/payment-service/internal/handler"
	"DewaSRY/go-microservices/services/payment-service/internal/repository"
	"DewaSRY/go-microservices/services/payment-service/internal/service"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
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
	GrpcAddr        = env.GetString("GRPC_ADDR", ":9004")
	rabbitMqURI     = env.GetString("AMQP_URL", "amqp://guest:guest@rabbitmq:5672/")
	appURL          = env.GetString("APP_URL", "http://localhost:3000")
	StripeSecretKey = env.GetString("STRIPE_SECRET_KEY", "")
	SuccessURL      = env.GetString("STRIPE_SUCCESS_URL", appURL+"?payment=success")
	CancelURL       = env.GetString("STRIPE_CANCEL_URL", appURL+"?payment=cancel")
	postgres_uri    = env.GetString("POSTGRES_URI", "postgres://postgres:postgres@postgres:5432/riderdb?sslmode=disable")
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

	// Stripe config
	stripeCfg := &types.PaymentConfig{
		StripeSecretKey: StripeSecretKey,
		SuccessURL:      SuccessURL,
		CancelURL:       CancelURL,
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

	db, err := db.NewPostgresManager(ctx, postgres_uri)
	if err != nil {
		log.Printf("failed_to_make_db_connection:%v", err)
		cancel()
		return
	}

	paymentRepo := repository.NewPaymentRepository(db)

	paymentProcessor := service.NewPaymentProcessService(stripeCfg)
	service := service.NewPaymentService(paymentRepo, paymentProcessor)
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
