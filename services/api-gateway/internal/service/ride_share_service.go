package service

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/services/api-gateway/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
)

type rideShareService struct {
	rabbitMq *messaging.RabbitMQ
}

// CreateTripEvent implements domain.RideShareServices.
func (t *rideShareService) CreateTripEvent(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.CreateTripRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		log.Printf("error_failed_to_process_user_init_data:%v", err)
		return nil
	}

	resultData, err := json.Marshal(
		messaging.CreateTripRequest{
			ConnectionId: connectionId,
			RiderId:      connectionId,
			Pickup:       parseData.Pickup,
			Destination:  parseData.Destination,
		},
	)

	if err != nil {
		log.Printf("failed_to_parse:%v", err)
		return nil
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.TripCreateInitProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		log.Printf("failed_to_parse:%v", err)
	}

	return nil
}

// UserInitEventRequest implements domain.RideShareServices.
func (t *rideShareService) UserInitEventRequest(ctx context.Context, connectionId string, data []byte) error {
	var parseData messaging.InitConnectionRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		log.Printf("error_failed_to_process_user_init_data:%v", err)
		return nil
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
		return nil
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.UserInitEventProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		log.Printf("failed_to_parse:%v", err)
	}

	return nil
}

func NewRideShareService(rabbitMq *messaging.RabbitMQ) domain.RideShareServices {
	return &rideShareService{rabbitMq: rabbitMq}
}
