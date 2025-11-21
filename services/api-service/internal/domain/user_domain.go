package domain

import (
	_type "DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"context"
)

type UserRepository interface {
	CreateRider(ctx context.Context, data _type.CreateRiderParam) error
	CreateDriver(ctx context.Context, data _type.CreateDriverParam) error
}

type UserEventHandler interface {
	HandlerUserInitConnection(ctx context.Context, data []byte)
	HandlerUserDisconnect(ctx context.Context, data []byte)
}

type UserService interface {
	UserInit(ctx context.Context, request messaging.InitConnectionRequest) error
}
