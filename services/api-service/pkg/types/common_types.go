package types

import sharedType "DewaSRY/go-microservices/shared/types"

type RoutesValue struct {
	Coordinate []sharedType.Coordinate
	Distance   float64
	Duration   float64
}
