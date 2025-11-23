package messaging

import "DewaSRY/go-microservices/shared/types"

type CommonSuccessResponse struct {
	Message string `json:"message"`
	Code    uint32 `json:"code"`
}

type RoutesResponse struct {
	Coordinate []types.Coordinate `json:"coordinate"`
	Distance   float64            `json:"distance"`
	Duration   float64            `json:"duration"`
}
