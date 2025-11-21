package service

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
)

type rideShareService struct {
	rabbitMq *messaging.RabbitMQ
}

// UserInitEvent implements domain.RideShareServices.
func (t *rideShareService) UserInitEvent(ctx context.Context, connectionId string, data []byte) {
	var parseData messaging.InitConnectionRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		log.Printf("error_failed_to_process_user_init_data:%v", err)
		return
	}

	resultData, err := json.Marshal(
		messaging.InitConnectionRequest{
			ConnectionId: connectionId,
			Coordinate:   parseData.Coordinate,
			PackageSlug:  parseData.PackageSlug,
			Entity:       parseData.Entity,
		},
	)

	if err != nil {
		log.Printf("failed_to_parse:%v", err)
		return
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.UserInitEvent,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		log.Printf("failed_to_parse:%v", err)
	}
}

func NewRideShareService(rabbitMq *messaging.RabbitMQ) domain.RideShareServices {
	return &rideShareService{rabbitMq: rabbitMq}
}
