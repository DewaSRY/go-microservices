package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
)

type tripFlowService struct {
	rabbitMq     *messaging.RabbitMQ
	tripFlowRepo domain.TripFlowRepository
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
		rabbitMq:     rabbitMq,
		tripFlowRepo: tripFlowRepo,
	}
}
