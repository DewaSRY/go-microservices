package handler

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
	"encoding/json"
	"log"
	"math/rand"
)

type tripEventHandler struct {
	service        domain.DriverService
	messageManager *messaging.RabbitMQ
}

func (t *tripEventHandler) FindSuitableDriver(ctx context.Context, data messaging.TripEventData) error {
	selectedSlug := data.Trip.SelectedFare.PackageSlug
	matchDriver := t.service.FindAvailableDrivers(selectedSlug)

	if len(matchDriver) == 0 {
		if err := t.messageManager.PublishingMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: data.Trip.UserID,
		}); err != nil {
			log.Printf("failed_to_publishes_message_to_exchange:%v", err)
			return err
		}

		return nil
	}

	tripDriverFindData, err := json.Marshal(
		messaging.TripDriverFindData{
			AmountDriver: len(matchDriver),
		},
	)
	if err != nil {
		return err
	}

	if err := t.messageManager.PublishingMessage(ctx, contracts.TripEventDriversFound,
		contracts.AmqpMessage{
			OwnerID: data.Trip.UserID,
			Data:    tripDriverFindData,
		},
	); err != nil {
		return err
	}

	// notify random driver
	randomIndex := rand.Intn(len(matchDriver))
	suitableDriverID := matchDriver[randomIndex]
	marshalledEvent, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := t.messageManager.PublishingMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: suitableDriverID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
		return err
	}

	return nil
}

func NewTripEventHandler(service domain.DriverService, messageManager *messaging.RabbitMQ) domain.DriverTripEventsHandler {
	return &tripEventHandler{service: service, messageManager: messageManager}
}
