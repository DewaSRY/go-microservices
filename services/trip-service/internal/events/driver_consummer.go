package events

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DriverConsumer struct {
	rabbitmq *messaging.RabbitMQ
	handler  domain.TripDriverEventHandler
}

func NewDriverConsumer(rabbitmq *messaging.RabbitMQ, handler domain.TripDriverEventHandler) *DriverConsumer {
	return &DriverConsumer{
		rabbitmq: rabbitmq,
		handler:  handler,
	}
}

func (t *DriverConsumer) Listen() error {
	return t.rabbitmq.ConsumeMessage(messaging.DriverTripResponseQueue, func(ctx context.Context, msg amqp.Delivery) error {
		var message contracts.AmqpMessage

		if err := json.Unmarshal(msg.Body, &message); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		switch msg.RoutingKey {
		case contracts.DriverCmdTripAccept:
			var payload messaging.DriverTripResponseData
			if err := json.Unmarshal(message.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			if err := t.handler.HandleTripAccepted(ctx, payload.TripID, payload.Driver); err != nil {
				log.Printf("failed_to_handle_the_trip_accept: %v", err)
				return err
			}
		case contracts.DriverCmdTripDecline:
			var payload messaging.DriverTripResponseData
			if err := json.Unmarshal(message.Data, &payload); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				return err
			}

			if err := t.handler.HandleTripDeclined(ctx, payload.TripID, payload.RiderId); err != nil {
				log.Printf("failed_to_handle_the_trip_decline: %v", err)
				return err
			}
		}

		return nil
	})
}
