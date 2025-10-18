package events

// NotifyPaymentSuccessQueue

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PaymentConsumer struct {
	rabbitmq *messaging.RabbitMQ
	handler  domain.TripPaymentEventHandler
}

func NewPaymentConsumer(rabbitmq *messaging.RabbitMQ, handler domain.TripPaymentEventHandler) *PaymentConsumer {
	return &PaymentConsumer{
		rabbitmq: rabbitmq,
		handler:  handler,
	}
}

func (t *PaymentConsumer) Listen() error {
	log.Print("get listen", messaging.NotifyPaymentSuccessQueue)
	return t.rabbitmq.ConsumeMessage(messaging.NotifyPaymentSuccessQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var message contracts.AmqpMessage

		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		// PaymentEventSuccess
		log.Print("get call", message)
		switch msg.RoutingKey {
		case contracts.PaymentEventSuccess:
			var payload messaging.AcceptedPayment
			if err := json.Unmarshal(message.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			log.Printf("p response received message: %v", payload)
			if err := t.handler.HandleAcceptedPayment(ctx, payload.TripID, payload.RiderID); err != nil {
				log.Printf("failed_to_handle_accepted_payment: %v", err)
				return err
			}
		}

		return nil
	})
}
