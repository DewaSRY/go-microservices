package domain

import (
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	"context"

	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"

	"gorm.io/gorm"
)

type DbFilter func(*gorm.DB) *gorm.DB

type DriverService interface {
	GetDriverProto(ctx context.Context, driverId string) (*drivergrpc.Driver, error)

	FindAvailableDrivers(ctx context.Context, packageTypes string) []string
	RegisterDriver(ctx context.Context, driverId string, packageSlug string) (*models.DriverModel, error)
	UnregisterDriver(ctx context.Context, driverId string) error
}

type DriverTripEventsHandler interface {
	FindSuitableDriver(context.Context, messaging.TripEventData) error
}

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *models.DriverModel) error
	GetActiveDriverIdList(ctx context.Context, filter DbFilter) ([]string, error)
	UpdateDriverData(ctx context.Context, driverId string, partialData map[string]interface{}) error

	GetDriverById(ctx context.Context, driverId string) (*models.DriverModel, error)
}
