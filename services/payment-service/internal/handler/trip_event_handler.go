package handler

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
)

type tripEventHandler struct {
	rabbitmq messaging.RabbitMQ
	service  domain.PaymentService
}

// HandleCreateSession implements domain.TripEventHandler.
func (t *tripEventHandler) HandleCreateSession(ctx context.Context, payload messaging.PaymentTripResponseData) error {
	paymentSession, err := t.service.CreatePaymentSession(
		ctx,
		payload.TripID,
		payload.UserID,
		payload.DriverID,
		int64(payload.Amount),
		payload.Currency,
	)
	if err != nil {
		log.Printf("Failed to create payment session: %v", err)
		return err
	}
	// Publish payment session created event
	paymentPayload := messaging.PaymentEventSessionCreatedData{
		TripID:    payload.TripID,
		SessionID: paymentSession.StripeSessionID,
		Amount:    float64(paymentSession.Amount) / 100.0, // Convert from cents to dollars
		Currency:  paymentSession.Currency,
	}

	payloadBytes, err := json.Marshal(paymentPayload)
	if err != nil {
		log.Printf("Failed to marshal payment session payload: %v", err)
		return err
	}

	if err := t.rabbitmq.PublishingMessage(
		ctx, contracts.PaymentEventSessionCreated,
		contracts.AmqpMessage{
			OwnerID: payload.UserID,
			Data:    payloadBytes,
		},
	); err != nil {
		log.Printf("Failed to publish payment session created event: %v", err)
		return err
	}

	log.Printf("Published payment session created event for trip: %s", payload.TripID)
	return nil
}

func NewTripEventHandler(
	rabbitmq messaging.RabbitMQ,
	service domain.PaymentService,
) domain.TripEventHandler {
	return &tripEventHandler{
		rabbitmq: rabbitmq,
		service:  service,
	}
}
