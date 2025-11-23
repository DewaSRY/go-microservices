package messaging

import "DewaSRY/go-microservices/shared/types"

type InitConnectionRequest struct {
	ConnectionId string           `json:"connectionId"`
	Coordinate   types.Coordinate `json:"coordinate"`
	PackageSlug  string           `json:"packageSlug"`
	Entity       string           `json:"entity"`
}

type CreateTripRequest struct {
	ConnectionId string           `json:"connectionId"`
	RiderId      string           `json:"riderId"`
	Pickup       types.Coordinate `json:"pickup"`
	Destination  types.Coordinate `json:"destination"`
}
