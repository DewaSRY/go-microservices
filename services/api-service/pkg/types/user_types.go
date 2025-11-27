package types

import "DewaSRY/go-microservices/shared/types"

type CreateRiderParam struct {
	Location     types.Coordinate
	Destination  types.Coordinate
	ConnectionId string
}

type CreateDriverParam struct {
	PackageSlug string
	Location    types.Coordinate
}

type UpdateRiderLocationParam struct {
	RiderId     string
	Location    types.Coordinate
	Destination types.Coordinate
}
