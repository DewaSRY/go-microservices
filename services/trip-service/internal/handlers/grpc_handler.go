package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
)

type grpcHandler struct {
	tripgrpc.UnimplementedTripServiceServer
	service domain.TripService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService) *grpcHandler {
	handler := &grpcHandler{service: service}

	tripgrpc.RegisterTripServiceServer(server, handler)
	return handler
}

func (t *grpcHandler) PreviewTrip(ctx context.Context, request *tripgrpc.PreviewTripRequest) (*tripgrpc.PreviewTripResponse, error) {
	log.Print("get call")

	route, err := t.service.GetRoute(ctx,
		&types.Coordinate{
			Latitude:  request.StartLocation.Latitude,
			Longitude: request.StartLocation.Longitude,
		}, &types.Coordinate{
			Latitude:  request.EndLocation.Latitude,
			Longitude: request.EndLocation.Longitude,
		})

	if err != nil {
		log.Printf("error %v", err)
		return nil, errors.New("get_error")
	}

	tripRideFare := &tripgrpc.RideFare{
		Id:                "",
		UserID:            request.UserID,
		PackageSlug:       "",
		TotalPriceInCents: 160,
	}

	currentDataRoute := route.Routes[0]
	responseCoordinate := make([]*tripgrpc.Coordinate, 0, len(currentDataRoute.Geometry.Coordinates))
	for _, currentCor := range currentDataRoute.Geometry.Coordinates {
		responseCoordinate = append(responseCoordinate, &tripgrpc.Coordinate{
			Latitude:  currentCor[0],
			Longitude: currentCor[1],
		})

	}

	responseRoute := &tripgrpc.Route{
		Geometry: &tripgrpc.Geometry{
			Coordinates: responseCoordinate,
		},
		Distance: currentDataRoute.Distance,
		Duration: currentDataRoute.Duration,
	}

	rideFareList := make([]*tripgrpc.RideFare, 0)
	rideFareList = append(rideFareList, tripRideFare)

	result := &tripgrpc.PreviewTripResponse{
		TripID:   "",
		RideFare: rideFareList,
		Route:    responseRoute,
	}

	return result, nil
}
