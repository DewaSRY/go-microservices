package events

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type tripConsume struct {
	messageManager *messaging.RabbitMQ
	handler        domain.DriverTripEventsHandler
}

func NewTripConsumer(messaging *messaging.RabbitMQ, handler domain.DriverTripEventsHandler) *tripConsume {
	return &tripConsume{messageManager: messaging, handler: handler}
}

func (t *tripConsume) Listen() error {
	return t.messageManager.ConsumeMessage(messaging.FindAvailableDriversQueue, func(ctx context.Context, d amqp.Delivery) error {
		var eventData contracts.AmqpMessage

		if err := json.NewDecoder(bytes.NewBuffer(d.Body)).Decode(&eventData); err != nil {
			return fmt.Errorf("failed_to_parse_event_message_data:%v", err)
		}

		switch d.RoutingKey {
		case contracts.TripEventCreated:
			var tripData messaging.TripEventData

			if err := json.NewDecoder(bytes.NewBuffer(eventData.Data)).Decode(&tripData); err != nil {
				return fmt.Errorf("failed_to_parse_trip_create_data:%v", err)
			}
			t.handler.FindSuitableDriver(ctx, tripData)
		case contracts.TripEventDriverNotInterested:
			var tripData messaging.TripEventData

			if err := json.NewDecoder(bytes.NewBuffer(eventData.Data)).Decode(&tripData); err != nil {
				return fmt.Errorf("failed_to_parse_trip_create_data:%v", err)
			}
			t.handler.FindSuitableDriver(ctx, tripData)
		}

		return nil
	})
}
