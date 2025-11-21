package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
)

type userService struct {
	rabbitMq *messaging.RabbitMQ
	userRepo domain.UserRepository
}

// UserInit implements domain.UserService.
func (t *userService) UserInit(ctx context.Context, request messaging.InitConnectionRequest) error {
	switch request.Entity {
	case "DRIVER":
		t.userRepo.CreateDriver(ctx, types.CreateDriverParam{
			PackageSlug: request.PackageSlug,
			Location:    request.Coordinate,
		})
	case "RIDER":
		t.userRepo.CreateRider(ctx, types.CreateRiderParam{
			PackageSlug: request.PackageSlug,
			Location:    request.Coordinate,
		})
	}

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
