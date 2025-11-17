package dto

import "DewaSRY/go-microservices/shared/types"

type UserInitRequest struct {
	PackageSlug string           `json:"packageSlug"`
	Location    types.Coordinate `json:"location"`
}
