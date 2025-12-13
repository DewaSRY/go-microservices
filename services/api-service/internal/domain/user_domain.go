package domain

import (
	_type "DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/models"
	"context"
)

type UserRepository interface {
	CreateRider(ctx context.Context, connectionId string, data _type.CreateRiderParam) error
	CreateDriver(ctx context.Context, connectionId string, data _type.CreateDriverParam) error
	UpdateRiderLocation(ctx context.Context, riderId string, location []byte, destination []byte) error

	CreateOrUpdateRiderModel(ctx context.Context, model models.RiderModel) error
	CreateOrUpdateDriverModel(ctx context.Context, model models.DriverModel) error

	GetDriverActiveList(ctx context.Context) ([]models.DriverModel, error)
	GetWaitingRiderIdConnectionList(ctx context.Context) ([]string, error)

	//TODO: clean up
	CleanUpDriverData(ctx context.Context, connectionId string) error
	CleanUpRiderData(ctx context.Context, connectionId string) error
}

type UserEventHandler interface {
	HandlerUserDisconnect(ctx context.Context, data []byte)
}

type UserService interface {
	UpdateRiderLocation(ctx context.Context, data _type.UpdateRiderLocationParam) error

	CreateRider(ctx context.Context, data _type.CreateRiderParam) error
	CreateDriver(ctx context.Context, data _type.CreateDriverParam) error

	NotifyDriverActive(ctx context.Context) error

	RiderStartTransaction(ctx context.Context, riderId string, driverId string) (string, error)
	DriverNotifyTransaction(ctx context.Context, driverId string, transactionId string) error
	RiderNotifyTransaction(ctx context.Context, connection string, transactionId string) error
	NotifyDriverAcceptedTransaction(ctx context.Context, transactionId string, driverId string, riderId string) error

	UserCleanUpData(ctx context.Context, connectionId string) error
}
