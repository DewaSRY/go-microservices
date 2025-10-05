package handlers

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	tripgrpc "DewaSRY/go-microservices/shared/proto/trip_proto"
	"DewaSRY/go-microservices/shared/types"
	"context"
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
		return nil, status.Errorf(codes.Internal, "failed_to_get_route:%v", err)
	}
	responseRoute := route.Routes[0].ToRouteProto()

	tripFareList := t.tripFareService.EstimatePackagesPrice(responseRoute.Distance, responseRoute.Duration)
	generatedTripFareList, err := t.tripFareService.GenerateTripFares(ctx, tripFareList, request.UserID, route)
	if err != nil {
		log.Printf("error %v", err)
		return nil, status.Errorf(codes.Internal, "failed_to_get_generate_ride_fare:%v", err)
	}

	rideFareList := make([]*tripgrpc.RideFare, len(generatedTripFareList))
	for i, fare := range generatedTripFareList {
		rideFareList[i] = &tripgrpc.RideFare{
			Id:                fare.ID.Hex(),
			UserID:            fare.UserID,
			PackageSlug:       fare.PackageSlug,
			TotalPriceInCents: fare.TotalPriceInCents,
		}
	}

	result := &tripgrpc.PreviewTripResponse{
		RideFares: rideFareList,
		Route:     responseRoute,
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
