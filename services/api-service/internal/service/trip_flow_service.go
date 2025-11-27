package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	_types "DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	"context"

	"github.com/google/uuid"
)

type tripFlowService struct {
	tripFlowRepo domain.TripFlowRepository
}

// CreateRiderTrip implements domain.TripFlowService.
func (t *tripFlowService) CreateRiderTrip(ctx context.Context, data _types.CreateTripParam) error {

	createModel := models.TripModel{
		RiderId: data.RiderId,
		Status:  "pending",
		Id:      uuid.New().String(),
	}

	return t.tripFlowRepo.CreateOrUpdateRiderTrip(ctx, createModel)
}

// TripCreate implements domain.TripFlowService.
func (t *tripFlowService) TripCreate(ctx context.Context, riderId string) error {
	return t.tripFlowRepo.CreateTrip(ctx, riderId)
}

func NewTripFlowService(
	rabbitMq *messaging.RabbitMQ,
	tripFlowRepo domain.TripFlowRepository,
) domain.TripFlowService {
	return &tripFlowService{
		tripFlowRepo: tripFlowRepo,
	}
}
