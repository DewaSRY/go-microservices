package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
)

type tripFlowService struct {
	rabbitMq     *messaging.RabbitMQ
	tripFlowRepo domain.TripFlowService
}

// TripCreate implements domain.TripFlowService.
func (t *tripFlowService) TripCreate(ctx context.Context, riderId string) error {
	return t.tripFlowRepo.TripCreate(ctx, riderId)
}

func NewtripFlowService(
	rabbitMq *messaging.RabbitMQ,
	tripFlowRepo domain.TripFlowService,
) domain.TripFlowService {
	return &tripFlowService{
		rabbitMq:     rabbitMq,
		tripFlowRepo: tripFlowRepo,
	}
}
