package domain

import (
	"DewaSRY/go-microservices/shared/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TripModel struct {
	ID       primitive.ObjectID
	UserID   string
	Status   string
	RideFare RideFareModel
}

type RideFareModel struct {
	ID                primitive.ObjectID
	UserID            string
	PackageSlug       string // ex : van, luxury. sedan
	TotalPriceInCents float64
	ExpiresAt         time.Time
	Route             *types.OsrmApiResponse
}

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *TripModel) (*TripModel, error)
	SaveRIdeFareList(ctx context.Context, fares []*RideFareModel) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *RideFareModel) (*TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
}

type TripFareService interface {
	EstimatePackagesPrice(distanceInKm float32, duration float64) []*RideFareModel
	GenerateTripFares(ctx context.Context, fares []*RideFareModel, userId string, route *types.OsrmApiResponse) ([]*RideFareModel, error)
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
