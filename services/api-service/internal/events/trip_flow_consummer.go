package events

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type TripFlowConsumer struct {
	rabbitmq        *messaging.RabbitMQ
	tripFlowHandler domain.TripFlowHandler
}

func (t *TripFlowConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessage(messaging.UserEstablishConnectionQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.MessageData
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.UserInitEvent:
			t.tripFlowHandler.HandlerTripCreate(ctx, message.Data)
		}

		return nil
	})
}

func NewTripFlowConsumer(rabbitmq *messaging.RabbitMQ, tripFlowHandler domain.TripFlowHandler) *TripFlowConsumer {
	return &TripFlowConsumer{rabbitmq: rabbitmq, tripFlowHandler: tripFlowHandler}
}
