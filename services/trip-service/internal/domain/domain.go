package domain

import (
	triptype "DewaSRY/go-microservices/services/trip-service/pkg/types"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
)

type TripRepository interface {
	CreateTrip(ctx context.Context, trip *triptype.TripModel) (*triptype.TripModel, error)
	SaveRIdeFareList(ctx context.Context, fares []*triptype.RideFareModel) error
	GetRideFareById(ctx context.Context, rideFareID string) (*triptype.RideFareModel, error)

	GetFareById(ctx context.Context, fareId string) (*triptype.RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*triptype.TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *triptype.RideFareModel) (*triptype.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
	GetUserRideFare(ctx context.Context, userID string, rideFareId string) (*triptype.RideFareModel, error)

	GetFareById(ctx context.Context, fareId string) (*triptype.RideFareModel, error)
	GetTripByID(ctx context.Context, id string) (*triptype.TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error
}

type TripFareService interface {
	EstimatePackagesPrice(distanceInKm float64, duration float64) []*triptype.RideFareModel
	GenerateTripFares(ctx context.Context, fares []*triptype.RideFareModel, userId string, route *types.OsrmApiResponse) ([]*triptype.RideFareModel, error)
	EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*triptype.RideFareModel
}

type TripDriverEventHandler interface {
	HandleTripAccepted(ctx context.Context, tripID string, driver *drivergrpc.Driver) error
	HandleTripDeclined(ctx context.Context, tripID string, riderId string) error
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
