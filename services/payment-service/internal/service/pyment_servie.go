package service

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type paymentService struct {
	paymentRepo      domain.PaymentRepo
	paymentProcessor domain.PaymentProcessorServiceService
}

// CreatePaymentSession implements domain.PaymentService.
func (t *paymentService) CreatePaymentSession(ctx context.Context, tripID string, userID string, driverID string, amount int64, currency string) (*types.PaymentIntent, error) {
	metadata := map[string]string{
		"trip_id":   tripID,
		"user_id":   userID,
		"driver_id": driverID,
	}

	sessionID, err := t.paymentProcessor.CreatePaymentSession(ctx, amount, currency, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	paymentIntent := &types.PaymentIntent{
		ID:              uuid.New().String(),
		TripID:          tripID,
		UserID:          userID,
		DriverID:        driverID,
		Amount:          amount,
		Currency:        currency,
		StripeSessionID: sessionID,
		CreatedAt:       time.Now(),
	}

	paymentModel := &models.TransactionModel{
		ID:       uuid.New().String(),
		RiderID:  userID,
		DriverID: driverID,
		Status:   "success",
	}

	if err := t.paymentRepo.CreateTransaction(ctx, paymentModel); err != nil {
		return nil, err
	}

	return paymentIntent, nil
}

func NewPaymentService(paymentRepo domain.PaymentRepo, paymentProcessor domain.PaymentProcessorServiceService) domain.PaymentService {
	return &paymentService{
		paymentRepo:      paymentRepo,
		paymentProcessor: paymentProcessor,
	}
}
