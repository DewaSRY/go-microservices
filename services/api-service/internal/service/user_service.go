package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"encoding/json"
)

type userService struct {
	rabbitMq *messaging.RabbitMQ
	userRepo domain.UserRepository
}

// CreateDriver implements domain.UserService.
func (t *userService) CreateDriver(ctx context.Context, data types.CreateDriverParam) error {

	locationJson, err := json.Marshal(data.Location)
	if err != nil {
		return err
	}

	newDriverModel := models.DriverModel{
		Id:          data.ConnectionId,
		Location:    locationJson,
		IsActive:    true,
		PackageSlug: data.PackageSlug,
	}

	return t.userRepo.CreateOrUpdateDriverModel(ctx, newDriverModel)

}

// CreateRider implements domain.UserService.
func (t *userService) CreateRider(ctx context.Context, data types.CreateRiderParam) error {

	locationJson, err := json.Marshal(data.Location)
	if err != nil {
		return err
	}

	destinationJson, err := json.Marshal(data.Destination)
	if err != nil {
		return err
	}

	newRiderModel := models.RiderModel{
		Id:          data.ConnectionId,
		Location:    locationJson,
		Destination: destinationJson,
		IsActive:    true,
	}

	return t.userRepo.CreateOrUpdateRiderModel(ctx, newRiderModel)
}

// UpdateRiderLocation implements domain.UserService.
func (t *userService) UpdateRiderLocation(ctx context.Context, data types.UpdateRiderLocationParam) error {

	byteLocation, err := json.Marshal(data.Location)
	if err != nil {
		return err
	}

	byteDestination, err := json.Marshal(data.Destination)
	if err != nil {
		return err
	}

	return t.userRepo.UpdateRiderLocation(ctx, data.RiderId, byteLocation, byteDestination)
}

// UserInit implements domain.UserService.
func (t *userService) UserInit(ctx context.Context, request messaging.InitConnectionRequest) error {
	// switch request.Entity {
	// case "DRIVER":
	// 	t.userRepo.CreateDriver(ctx, request.ConnectionId, types.CreateDriverParam{
	// 		PackageSlug: request.PackageSlug,
	// 		Location:    request.Coordinate,
	// 	})
	// case "RIDER":
	// 	t.userRepo.CreateRider(ctx, request.ConnectionId, types.CreateRiderParam{
	// 		PackageSlug: request.PackageSlug,
	// 		Location:    request.Coordinate,
	// 	})
	// }

	return nil
}

func NewUserService(
	rabbitMq *messaging.RabbitMQ,
	userRepo domain.UserRepository,
) domain.UserService {
	return &userService{
		rabbitMq: rabbitMq,
		userRepo: userRepo,
	}
}
