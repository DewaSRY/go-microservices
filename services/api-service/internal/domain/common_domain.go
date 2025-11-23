package domain

import (
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	sharedType "DewaSRY/go-microservices/shared/types"
	"context"
)

type OsrmIntegrationService interface {
	GetRoutes(ctx context.Context, start sharedType.Coordinate, end sharedType.Coordinate) (*types.RoutesValue, error)
}
