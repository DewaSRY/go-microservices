package service

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/contracts"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	_types "DewaSRY/go-microservices/shared/types"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type userService struct {
	rabbitMq *messaging.RabbitMQ
	userRepo domain.UserRepository
	tripRepo domain.TripFlowRepository
}

// UserCleanUpData implements domain.UserService.
func (t *userService) UserCleanUpData(ctx context.Context, connectionId string) error {

	if err := t.userRepo.CleanUpDriverData(ctx, connectionId); err != nil {
		return err
	}

	if err := t.userRepo.CleanUpRiderData(ctx, connectionId); err != nil {
		return nil
	}

	if err := t.tripRepo.CleanUpTransaction(ctx, connectionId); err != nil {
		return nil
	}

	return nil
}

// NotifyDriverAcceptedTransaction implements domain.userService.
func (t *userService) NotifyDriverAcceptedTransaction(ctx context.Context, transactionId string, driverId string, riderId string) error {
	responseJson, err := json.Marshal(messaging.TransactionAcceptedResponse{
		TransactionId: transactionId,
	})

	if err != nil {
		return err
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.TransactionAcceptedResponse, contracts.MessageData{
		ConnectionId: driverId,
		Data:         responseJson,
	}); err != nil {
		return err
	}

	if err := t.rabbitMq.PublishingMessage(ctx, contracts.TransactionAcceptedResponse, contracts.MessageData{
		ConnectionId: riderId,
		Data:         responseJson,
	}); err != nil {
		return err
	}

	return nil
}

// RiderNotifyTransaction implements domain.UserService.
func (t *userService) RiderNotifyTransaction(ctx context.Context, connection string, transactionId string) error {

	jsonData, err := json.Marshal(messaging.RiderCreateTransactionResponse{
		TransactionId: transactionId,
	})

	if err != nil {
		return nil
	}

	if err := t.rabbitMq.PublishingMessage(
		ctx, contracts.RiderCreateTransactionResponse,
		contracts.MessageData{
			ConnectionId: connection,
			Data:         jsonData,
		}); err != nil {
		return nil
	}
	return nil
}

// DriverNotifyTransaction implements domain.UserService.
func (t *userService) DriverNotifyTransaction(ctx context.Context, driverId string, transactionId string) error {

	jsonData, err := json.Marshal(messaging.RiderCreateTransactionResponse{
		TransactionId: transactionId,
	})

	if err != nil {
		return nil
	}

	if err := t.rabbitMq.PublishingMessage(
		ctx, contracts.RiderCreateTransactionResponse,
		contracts.MessageData{
			ConnectionId: driverId,
			Data:         jsonData,
		}); err != nil {
		return nil
	}
	return nil
}

// RiderStartTransaction implements domain.UserService.
func (t *userService) RiderStartTransaction(ctx context.Context, riderId string, driverId string) (string, error) {

	transactionId := uuid.New().String()
	transactionModel := models.TransactionModel{
		Id:       transactionId,
		RiderId:  riderId,
		DriverId: driverId,
		Status:   "PENDING",
	}

	if err := t.tripRepo.CreateOrUpdateTransactionModel(ctx, transactionModel); err != nil {
		return "", err
	}

	return transactionId, nil
}

// NotifyDriverActive implements domain.UserService.
func (t *userService) NotifyDriverActive(ctx context.Context) error {
	driverList, err := t.userRepo.GetDriverActiveList(ctx)

	if err != nil {
		return err
	}

	waitingRiderList, err := t.userRepo.GetWaitingRiderIdConnectionList(ctx)
	if err != nil {
		return err
	}

	if len(waitingRiderList) == 0 {
		return nil
	}

	driverRecordList := make([]messaging.DriverRecordResponse, 0)

	for _, d := range driverList {
		var unMarshelLocation _types.Coordinate

		if err := json.Unmarshal(d.Location, &unMarshelLocation); err != nil {
			return nil
		}

		driverRecordList = append(driverRecordList, messaging.DriverRecordResponse{
			Coordinate:  unMarshelLocation,
			PackageSlug: d.PackageSlug,
			DriverId:    d.Id,
		})
	}

	jsonData, err := json.Marshal(driverRecordList)

	if err != nil {
		return err
	}

	for _, ids := range waitingRiderList {
		if err := t.rabbitMq.PublishingMessage(ctx, contracts.DriverActiveResponse, contracts.MessageData{
			ConnectionId: ids,
			Data:         jsonData,
		}); err != nil {
			return err
		}
	}
	return nil
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

func NewUserService(
	rabbitMq *messaging.RabbitMQ,
	userRepo domain.UserRepository,
	tripRepo domain.TripFlowRepository,
) domain.UserService {
	return &userService{
		rabbitMq: rabbitMq,
		userRepo: userRepo,
		tripRepo: tripRepo,
	}
}
