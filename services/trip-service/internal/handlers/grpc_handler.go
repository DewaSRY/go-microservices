package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	tripgrpc.UnimplementedTripServiceServer
	service         domain.TripService
	tripFareService domain.TripFareService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService, tripFareService domain.TripFareService) *grpcHandler {
	handler := &grpcHandler{service: service, tripFareService: tripFareService}

	tripgrpc.RegisterTripServiceServer(server, handler)
	return handler
}

func (t *grpcHandler) CreateTrip(ctx context.Context, request *tripgrpc.CreateTripRequest) (*tripgrpc.CreateTripResponse, error) {
	createdTrip, err := t.service.CreateTrip(ctx, &domain.RideFareModel{
		ID:                primitive.NewObjectID(),
		UserID:            request.UserID,
		PackageSlug:       "some-test",
		TotalPriceInCents: 18,
		ExpiresAt:         time.Now(),
	})

	if err != nil {
		log.Printf("error %v", err)
		return nil, errors.New("get_error")
	}

	response := &tripgrpc.CreateTripResponse{
		TripID: createdTrip.ID.Hex(),
	}
	return response, nil
}

func (t *grpcHandler) PreviewTrip(ctx context.Context, request *tripgrpc.PreviewTripRequest) (*tripgrpc.PreviewTripResponse, error) {
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

	// 1. Estimation the ride fare price base on the route (ex distance)
	// 2. store the ride fares for the create trip (next leason) to fetch and validation
	tripFareList := t.tripFareService.EstimatePackagesPrice(
		float32(currentDataRoute.Distance),
		currentDataRoute.Duration,
	)

	generatedTripFareList, err := t.tripFareService.GenerateTripFares(ctx, tripFareList, request.UserID, route)
	if err != nil {
		log.Printf("error %v", err)
		return nil, errors.New("get_error")
	}

	rideFareList := make([]*tripgrpc.RideFare, 0, len(generatedTripFareList))

	for i, fare := range generatedTripFareList {
		rideFareList[i] = &tripgrpc.RideFare{
			Id:                fare.ID.Hex(),
			UserID:            fare.UserID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}
	}

	result := &tripgrpc.PreviewTripResponse{
		TripID:   "",
		RideFare: rideFareList,
		Route:    responseRoute,
	}

	return result, nil
}

package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	tripgrpc.UnimplementedTripServiceServer
	service         domain.TripService
	tripFareService domain.TripFareService
}

func NewGRPCHandler(server *grpc.Server, service domain.TripService, tripFareService domain.TripFareService) *grpcHandler {
	handler := &grpcHandler{service: service, tripFareService: tripFareService}

	tripgrpc.RegisterTripServiceServer(server, handler)
	return handler
}

func (t *grpcHandler) CreateTrip(ctx context.Context, request *tripgrpc.CreateTripRequest) (*tripgrpc.CreateTripResponse, error) {
	userID := request.GetUserID()
	rideFareID := request.GetRideFareID()

	userRideFare, err := t.service.GetUserRideFare(ctx, userID, rideFareID)
	if err != nil {
		log.Printf("error %v", err)
		return nil, status.Errorf(codes.Internal, "failed_to_get_user_route:%v", err)
	}

	createdTrip, err := t.service.CreateTrip(ctx, userRideFare)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed_to_create_user_error: %v", err)
	}

	response := &tripgrpc.CreateTripResponse{
		TripID: createdTrip.ID.Hex(),
	}
	return response, nil
}

func (t *grpcHandler) PreviewTrip(ctx context.Context, request *tripgrpc.PreviewTripRequest) (*tripgrpc.PreviewTripResponse, error) {
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

	responseRoute := route.Routes[0].ToRouteProto()
	// 1. Estimation the ride fare price base on the route (ex distance)
	// 2. store the ride fares for the create trip (next leason) to fetch and validation
	tripFareList := t.tripFareService.EstimatePackagesPrice(responseRoute.Distance, responseRoute.Duration)
	generatedTripFareList, err := t.tripFareService.GenerateTripFares(ctx, tripFareList, request.UserID, route)
	if err != nil {
		log.Printf("error %v", err)
		return nil, errors.New("get_error")
	}

	rideFareList := make([]*tripgrpc.RideFare, 0, len(generatedTripFareList))

	for i, fare := range generatedTripFareList {
		rideFareList[i] = &tripgrpc.RideFare{
			Id:                fare.ID.Hex(),
			UserID:            fare.UserID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}
	}

	result := &tripgrpc.PreviewTripResponse{
		TripID:   "",
		RideFare: rideFareList,
		Route:    responseRoute,
	}

	return result, nil
}
