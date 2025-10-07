package domain

import "DewaSRY/go-microservices/services/driver-service/pkg"

type DriverService interface {
	FindAvailableDrivers(packageTypes string) []string
	RegisterDriver(driverId string, packageSlug string) (pkg.Driver, error)
	UnregisterDriver(driverId string)
}

type DriverRepository interface {
}
