package events

import (
	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
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

func (p *TripEventPublisher) PublishTripCreated(ctx context.Context, tripModel *triptype.TripModel) error {
	payload, err := json.Marshal(
		messaging.TripEventData{
			Trip: tripModel.ToTripProto(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed_to_publish_trip_create:%v", err)
	}
	return p.messageManager.PublishingMessage(ctx, contracts.TripEventCreated,
		contracts.AmqpMessage{
			OwnerID: tripModel.UserID,
			Data:    payload,
		},
	)
}
