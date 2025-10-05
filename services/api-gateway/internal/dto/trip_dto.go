package dto

import "DewaSRY/go-microservices/shared/types"

type PreviewTripRequest struct {
	UserID      string           `json:"userID" validate:"required"`
	Pickup      types.Coordinate `json:"pickup" validate:"required"`
	Destination types.Coordinate `json:"destination" validate:"required"`
}

type StartTripRequest struct {
	RideFareID string `json:"rideFareID" validate:"required"`
	UserID     string `json:"userID" validate:"required"`
}
