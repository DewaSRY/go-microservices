package types

import "DewaSRY/go-microservices/shared/types"

type UserInitRequest struct {
	PackageSlug string           `json:"packageSlug"`
	Location    types.Coordinate `json:"location"`
}

type CreateTripRequest struct {
	RiderId     string           `json:"riderId"`
	Pickup      types.Coordinate `json:"pickup"`
	Destination types.Coordinate `json:"destination"`
}
