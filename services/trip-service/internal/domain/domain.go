package domain

import (
	"DewaSRY/go-microservices/shared/models"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"

	"gorm.io/gorm"
)

type DbFilter func(*gorm.DB) *gorm.DB

type TripRepository interface {
	FindTripWithFilter(dbFilter DbFilter) (*models.TripModel, error)
	GetFareById(ctx context.Context, fareId string) (*models.FareModel, error)
	GetTripByID(ctx context.Context, id string) (*models.TripModel, error)

	CreateTrip(ctx context.Context, trip *models.TripModel) error
	CreateFareList(ctx context.Context, fares []*models.FareModel) error
	CreateFare(ctx context.Context, fare *models.FareModel) error

	UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error
}

type TripService interface {
	CreateTrip(ctx context.Context, fare *models.FareModel) (*models.TripModel, error)
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate) (*types.OsrmApiResponse, error)
	GetUserRideFare(ctx context.Context, userID string, rideFareId string) (*models.FareModel, error)
	GetUserTrip(ctx context.Context, userId string, fareId string) (*models.TripModel, error)

	GetFareById(ctx context.Context, fareId string) (*models.FareModel, error)
	GetTripByID(ctx context.Context, id string) (*models.TripModel, error)
	UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error

	GetTripProto(ctx context.Context, tripId string) (*tripgrpc.Trip, error)
}

type TripFareService interface {
	EstimatePackagesPrice(distanceInKm float64, duration float64) []*models.FareModel
	GenerateTripFares(ctx context.Context, fares []*models.FareModel, userId string, route *types.OsrmApiResponse) ([]*models.FareModel, error)
	EstimatePackagesPriceWithRoute(route *types.OsrmApiResponse) []*models.FareModel
}

type TripDriverEventHandler interface {
	HandleTripAccepted(ctx context.Context, tripID string, driver *drivergrpc.Driver) error
	HandleTripDeclined(ctx context.Context, tripID string, riderId string) error
}

type TripPaymentEventHandler interface {
	HandleAcceptedPayment(ctx context.Context, TripID, RiderID string) error
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
