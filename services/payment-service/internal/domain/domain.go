package domain

import (
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"DewaSRY/go-microservices/shared/messaging"
	"DewaSRY/go-microservices/shared/models"
	"context"
)

type PaymentRepo interface {
	CreateTransaction(ctx context.Context, transaction *models.TransactionModel) error
}

type PaymentService interface {
	CreatePaymentSession(ctx context.Context, tripID, userID, driverID string, amount int64, currency string) (*types.PaymentIntent, error)
}

type PaymentProcessorServiceService interface {
	CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error)
}

type TripEventHandler interface {
	HandleCreateSession(ctx context.Context, payload messaging.PaymentTripResponseData) error
}
