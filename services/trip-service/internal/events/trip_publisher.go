package events

import (
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"context"
	"encoding/json"
	"fmt"
)

type TripEventPublisher struct {
	messageManager *messaging.RabbitMQ
}

func NewTripEventPublisher(messageManager *messaging.RabbitMQ) *TripEventPublisher {
	return &TripEventPublisher{messageManager: messageManager}
}

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, Trip *tripgrpc.Trip) error {
	payload, err := json.Marshal(
		messaging.TripEventData{
			Trip: Trip,
		},
	)
	if err != nil {
		return fmt.Errorf("failed_to_publish_trip_create:%v", err)
	}
	return p.messageManager.PublishingMessage(ctx, contracts.TripEventCreated,
		contracts.AmqpMessage{
			OwnerID: Trip.UserID,
			Data:    payload,
		},
	)
}
