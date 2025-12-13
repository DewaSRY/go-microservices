package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/integrations"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	sharedType "DewaSRY/go-microservices/shared/types"
	"context"
)

type osrmIntegrationService struct {
}

// GetRoutes implements domain.OsrmIntegrationService.
func (o *osrmIntegrationService) GetRoutes(ctx context.Context, start sharedType.Coordinate, end sharedType.Coordinate) (*types.RoutesValue, error) {
	osrmRoute, err := integrations.GetRoute(ctx, &start, &end)
	if err != nil {
		return nil, err
	}

	coordinateList := make([]sharedType.Coordinate, 0)
	optimizeRoute := osrmRoute.Routes[0]

	for _, coord := range optimizeRoute.Geometry.Coordinates {
		coordinateList = append(coordinateList, sharedType.Coordinate{
			Latitude:  coord[1],
			Longitude: coord[0],
		})
	}

	return &types.RoutesValue{
		Coordinate: coordinateList,
		Distance:   optimizeRoute.Distance,
		Duration:   optimizeRoute.Duration,
	}, nil
}

func NewOsrmIntegrationService() domain.OsrmIntegrationService {
	return &osrmIntegrationService{}
}
