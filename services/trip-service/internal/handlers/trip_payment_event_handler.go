package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"fmt"
)

type tripPaymentEventHandler struct {
	rabbitmq    *messaging.RabbitMQ
	tripService domain.TripService
}

// HandleAcceptedPayment implements domain.TripPaymentEventHandler.
func (t *tripPaymentEventHandler) HandleAcceptedPayment(ctx context.Context, TripID, RiderID string) error {
	// panic("unimplemented")

	trip, err := t.tripService.GetTripByID(ctx, TripID)

	if err != nil {
		return fmt.Errorf("trip_with_id_%s_not_found:%v", TripID, err)
	}

	driverId := trip.DriverId

	// contracts.PaymentEventComplete

	// temp pyload
	dataPayload, err := json.Marshal(map[string]string{
		"riderId": RiderID,
		"TripID":  TripID,
	})

	if err != nil {
		return fmt.Errorf("failed_to_marshal_payload:%v", err)
	}

	// notify driver
	if err := t.rabbitmq.PublishingMessage(
		ctx,
		contracts.PaymentEventComplete,
		contracts.AmqpMessage{
			OwnerID: driverId,
			Data:    dataPayload,
		},
	); err != nil {
		return err
	}

	// notify rider
	if err := t.rabbitmq.PublishingMessage(
		ctx,
		contracts.PaymentEventComplete,
		contracts.AmqpMessage{
			OwnerID: RiderID,
			Data:    dataPayload,
		},
	); err != nil {
		return err
	}

	return nil
}

func NewTripPaymentEventHandler(
	rabbitmq *messaging.RabbitMQ,
	tripService domain.TripService,
) domain.TripPaymentEventHandler {
	return &tripPaymentEventHandler{
		rabbitmq:    rabbitmq,
		tripService: tripService,
	}
}
