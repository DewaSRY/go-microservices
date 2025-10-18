package domain

import (
	"DewaSRY/go-microservices/services/driver-service/pkg"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
)

type DriverService interface {
	FindAvailableDrivers(packageTypes string) []string
	RegisterDriver(driverId string, packageSlug string) (pkg.Driver, error)
	UnregisterDriver(driverId string)
}

type DriverTripEventsHandler interface {
	FindSuitableDriver(context.Context, messaging.TripEventData) error
}

type DriverRepository interface {
}
