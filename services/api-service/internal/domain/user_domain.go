package domain

import (
	_type "DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	"context"
)

type UserRepository interface {
	CreateRider(ctx context.Context, connectionId string, data _type.CreateRiderParam) error
	CreateDriver(ctx context.Context, connectionId string, data _type.CreateDriverParam) error
	UpdateRiderLocation(ctx context.Context, riderId string, location []byte, destination []byte) error

	CreateOrUpdateRiderModel(ctx context.Context, model models.RiderModel) error
}

type UserEventHandler interface {
	HandlerUserInitConnection(ctx context.Context, data []byte)
	HandlerUserDisconnect(ctx context.Context, data []byte)
}

type UserService interface {
	UserInit(ctx context.Context, request messaging.InitConnectionRequest) error
	UpdateRiderLocation(ctx context.Context, data _type.UpdateRiderLocationParam) error

	CreateRider(ctx context.Context, data _type.CreateRiderParam) error
}
