package events

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type TripConsumer struct {
	rabbitmq messaging.RabbitMQ
	handler  domain.TripEventHandler
}

func (t *TripConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessage(messaging.PaymentTripResponseQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.PaymentTripResponseData
		if err := json.Unmarshal(message.Data, &payload); err != nil {
			log.Printf("Failed to unmarshal payload: %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.PaymentCmdCreateSession:
			if err := t.handler.HandleCreateSession(ctx, payload); err != nil {
				log.Printf("Failed to handle trip accepted: %v", err)
				return err
			}
		}

		return nil
	})
}

func NewTripConsumer(rabbitmq messaging.RabbitMQ, handler domain.TripEventHandler) *TripConsumer {
	return &TripConsumer{
		rabbitmq: rabbitmq,
		handler:  handler,
	}
}
