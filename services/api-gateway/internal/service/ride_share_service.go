package service

import (
	"DewaSRY/go-microservices/services/api-gateway/internal/domain"
	"DewaSRY/go-microservices/services/api-gateway/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"errors"
)

type rideShareService struct {
	rabbitMq *messaging.RabbitMQ
}

// UserDisconnected implements domain.RideShareServices.
func (t *rideShareService) UserDisconnected(ctx context.Context, connectionId string) error {
	messageData, err := json.Marshal(
		messaging.UserDisconnectionRequest{
			ConnectionId: connectionId,
		},
	)

	if err != nil {
		return errors.New("failed_to_parse_payload")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.UserDisconnectedProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         messageData,
		}); err != nil {
		return errors.New("failed_to_publish_message")
	}

	return nil
}

// DriverAcceptedTransaction implements domain.RideShareServices.
func (t *rideShareService) DriverAcceptedTransaction(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.DriverAcceptedTransactionRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		return errors.New("error_failed_to_process_user_init_data")
	}

	resultData, err := json.Marshal(
		messaging.DriverAcceptedTransactionRequest{
			ConnectionId:  connectionId,
			TransactionId: parseData.TransactionId,
		},
	)

	if err != nil {
		return errors.New("failed_to_parse")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.DriverAcceptTransactionProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		return errors.New("failed_to_send_message")
	}

	return nil
}

// RiderCreateTransaction implements domain.RideShareServices.
func (t *rideShareService) RiderCreateTransaction(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.RiderCreateTransactionRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		return errors.New("error_failed_to_process_user_init_data")
	}

	resultData, err := json.Marshal(
		messaging.RiderCreateTransactionRequest{
			ConnectionId: connectionId,
			DriverId:     parseData.DriverId,
		},
	)

	if err != nil {
		return errors.New("failed_to_parse")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.RiderCreateTransactionProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		return errors.New("failed_to_send_message")
	}

	return nil
}

// DriverInitRequest implements domain.RideShareServices.
func (t *rideShareService) DriverInitRequest(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.DriverInitRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		return errors.New("failed_to_read_data")
	}

	resultData, err := json.Marshal(
		messaging.DriverInitRequest{
			ConnectionId: connectionId,
			Location:     parseData.Location,
			PackageSlug:  parseData.PackageSlug,
		},
	)

	if err != nil {
		return errors.New("failed_to_parse")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.DriverInitEventProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		return errors.New("failed_to_send_message")
	}

	return nil
}

// RiderCreateTripRequest implements domain.RideShareServices.
func (t *rideShareService) RiderCreateTripRequest(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.RiderCreateTripRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		return errors.New("error_failed_to_process_user_init_data")
	}

	resultData, err := json.Marshal(
		messaging.RiderCreateTripRequest{
			ConnectionId: connectionId,
			Pickup:       parseData.Pickup,
			Destination:  parseData.Destination,
		},
	)

	if err != nil {
		return errors.New("failed_to_parse")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.RiderCreateTripProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		return errors.New("failed_to_send_message")
	}

	return nil
}

// CreateTripEvent implements domain.RideShareServices.
func (t *rideShareService) CreateTripEvent(ctx context.Context, connectionId string, data []byte) error {
	var parseData types.CreateTripRequest

	if err := json.Unmarshal(data, &parseData); err != nil {
		return errors.New("error_failed_to_process_user_init_data")
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
		return errors.New("failed_to_parse")
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.TripCreateInitProcess,
		contracts.MessageData{
			ConnectionId: connectionId,
			Data:         resultData,
		}); err != nil {
		return errors.New("failed_to_send_message")
	}

	return nil
}

func NewRideShareService(rabbitMq *messaging.RabbitMQ) domain.RideShareServices {
	return &rideShareService{rabbitMq: rabbitMq}
}
