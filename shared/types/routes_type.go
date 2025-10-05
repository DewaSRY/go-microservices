package types

import (
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
)

type Routes struct {
	Legs       []Legs   `json:"legs"`
	WeightName string   `json:"weight_name"`
	Geometry   Geometry `json:"geometry"`
	Weight     float64  `json:"weight"`
	Duration   float64  `json:"duration"`
	Distance   float64  `json:"distance"`
}

func (r *Routes) ToRouteProto() *tripgrpc.Route {
	responseCoordinate := make([]*tripgrpc.Coordinate, 0, len(r.Geometry.Coordinates))
	for _, currentCor := range r.Geometry.Coordinates {
		responseCoordinate = append(responseCoordinate, &tripgrpc.Coordinate{
			Latitude:  currentCor[0],
			Longitude: currentCor[1],
		})

	}

	return &tripgrpc.Route{
		Geometry: &tripgrpc.Geometry{
			Coordinates: responseCoordinate,
		},
		Distance: r.Distance,
		Duration: r.Duration,
	}
}
