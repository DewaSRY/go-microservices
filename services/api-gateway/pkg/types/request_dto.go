package types

import "DewaSRY/go-microservices/shared/types"

type UserInitRequest struct {
	PackageSlug string           `json:"packageSlug"`
	Location    types.Coordinate `json:"location"`
}

type CreateTripRequest struct {
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type RiderCreateTripRequest struct {
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}

type DriverInitRequest struct {
	Location    types.Coordinate `json:"location"`
	PackageSlug string           `json:"packageSlug"`
}
