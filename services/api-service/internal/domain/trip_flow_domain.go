package domain

import (
	"context"
)

type TripFlowRepository interface {
	CreateTrip(ctx context.Context, riderId string) error
}

type TripFlowHandler interface {
	HandlerTripCreate(ctx context.Context, data []byte)
}

type TripFlowService interface {
	TripCreate(ctx context.Context, riderId string) error
}
