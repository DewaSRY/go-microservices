package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"context"
	"encoding/json"
	"fmt"
)

type tripDriverEventHandler struct {
	rabbitmq    *messaging.RabbitMQ
	tripService domain.TripService
	fareService domain.TripFareService
}

// HandleTripAccepted implements domain.TripDriverEventHandler.
func (t *tripDriverEventHandler) HandleTripAccepted(ctx context.Context, tripID string, driver *drivergrpc.Driver) error {
	// trip, err := t.service.GenerateTripFares()
	trip, err := t.tripService.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed_to_get_trip_with_id_%s:%v", tripID, err)
	}

	if trip == nil {
		return fmt.Errorf("trip_with_id%s_not_found", tripID)
	}

	if err := t.tripService.UpdateTrip(ctx, tripID, "accepted", driver); err != nil {
		return fmt.Errorf("failed_to_update_the_trip:%v", err)
	}

	trip, err = t.tripService.GetTripByID(ctx, tripID)
	if err != nil {
		return fmt.Errorf("failed_to_get_trip_with_id_%s:%v", tripID, err)
	}

	marshalledTrip, err := json.Marshal(trip)
	if err != nil {
		return err
	}

	//Notify the rider the a driver has been assigned
	if err := t.rabbitmq.PublishingMessage(
		ctx,
		contracts.TripEventDriverAssigned,
		contracts.AmqpMessage{
			OwnerID: trip.UserID,
			Data:    marshalledTrip,
		}); err != nil {
		return err
	}
	marshalledPayload, err := json.Marshal(
		messaging.PaymentTripResponseData{
			TripID:   tripID,
			UserID:   trip.UserID,
			DriverID: driver.Id,
			Amount:   trip.RideFare.TotalPriceInCents,
			Currency: "USD",
		},
	)
	if err != nil {
		return fmt.Errorf("failed_to_marshal_payload:%v", err)
	}

	if err := t.rabbitmq.PublishingMessage(
		ctx,
		contracts.PaymentCmdCreateSession,
		contracts.AmqpMessage{
			OwnerID: trip.UserID,
			Data:    marshalledPayload,
		},
	); err != nil {
		return err
	}

	return nil
}

// HandleTripDeclined implements domain.TripDriverEventHandler.
func (t *tripDriverEventHandler) HandleTripDeclined(ctx context.Context, tripID string, riderId string) error {
	trip, err := t.tripService.GetTripByID(ctx, tripID)

	if err != nil {
		return fmt.Errorf("failed_to_get_trip_with_id_%s:%v", tripID, err)
	}

	newPayload := messaging.TripEventData{
		Trip: trip.ToTripProto(),
	}

	marshalledPayload, err := json.Marshal(newPayload)
	if err != nil {
		return fmt.Errorf("failed_to_marshal_with_error:%v", err)
	}

	if err := t.rabbitmq.PublishingMessage(
		ctx,
		contracts.TripEventDriverNotInterested,
		contracts.AmqpMessage{
			OwnerID: riderId,
			Data:    marshalledPayload,
		}); err != nil {
		return err
	}
	return nil

}

func NewTripDriverEventHandler(
	rabbitmq *messaging.RabbitMQ,
	tripService domain.TripService,
	fareService domain.TripFareService,
) domain.TripDriverEventHandler {
	return &tripDriverEventHandler{
		rabbitmq:    rabbitmq,
		tripService: tripService,
		fareService: fareService,
	}
}
