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

type UserConnectionConsumer struct {
	rabbitmq    *messaging.RabbitMQ
	userHandler domain.UserEventHandler
}

func (t *UserConnectionConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessage(messaging.UserEstablishConnectionQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var message contracts.MessageData
		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		log.Printf("test_:%v", msg.RoutingKey)

		switch msg.RoutingKey {
		case contracts.UserInitEventProcess:
			t.userHandler.HandlerUserInitConnection(ctx, message.Data)
		}

		return nil
	})
}

func NewUserConsumer(rabbitmq *messaging.RabbitMQ, userHandler domain.UserEventHandler) *UserConnectionConsumer {
	return &UserConnectionConsumer{rabbitmq: rabbitmq, userHandler: userHandler}
}
