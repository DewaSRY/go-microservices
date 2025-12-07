package domain

import (
	_types "DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/models"
	"context"
)

type TripFlowRepository interface {
	CreateTrip(ctx context.Context, riderId string) error

	CreateOrUpdateRiderTrip(ctx context.Context, model models.TripModel) error

	CreateOrUpdateTransactionModel(ctx context.Context, model models.TransactionModel) error
}

type TripFlowHandler interface {
	HandlerTripCreate(ctx context.Context, data []byte)
	HandlerRiderCreateTrip(ctx context.Context, data []byte)
	HandlerDriverInit(ctx context.Context, data []byte)
	HandleRiderCreateTransaction(ctx context.Context, data []byte)
}

type TripFlowService interface {
	TripCreate(ctx context.Context, riderId string) error
	CreateRiderTrip(ctx context.Context, data _types.CreateTripParam) error
}
