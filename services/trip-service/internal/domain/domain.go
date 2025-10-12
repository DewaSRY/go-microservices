package domain

import (
	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"
	"DewaSRY/go-microservices/shared/types"
	"context"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *triptype.TripModel) (*triptype.TripModel, error)
	SaveRIdeFareList(ctx context.Context, fares []*triptype.RideFareModel) error
	GetRideFareById(ctx context.Context, rideFareID string) (*triptype.RideFareModel, error)
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *triptype.RideFareModel) (*triptype.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
	GetUserRideFare(ctx context.Context, userID string, rideFareId string) (*triptype.RideFareModel, error)
}

type TripFareService interface {
	EstimatePackagesPrice(distanceInKm float64, duration float64) []*triptype.RideFareModel
	GenerateTripFares(ctx context.Context, fares []*triptype.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*triptype.RideFareModel, error)
}

type PricingConfig struct {
	PricePerUnitOfDistance float64
	PricingPerMinute       float64
}

func DefaultPricingConfig() *PricingConfig {
	return &PricingConfig{
		PricePerUnitOfDistance: 1.5,
		PricingPerMinute:       0.25,
	}
}
