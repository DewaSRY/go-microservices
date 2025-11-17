package types

import "DewaSRY/go-microservices/shared/types"

type CreateRiderParam struct {
	PackageSlug string
	Location    types.Coordinate
}

type CreateDriverParam struct {
	PackageSlug string
	Location    types.Coordinate
}
