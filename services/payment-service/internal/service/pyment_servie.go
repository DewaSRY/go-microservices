package service

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type paymentService struct {
	paymentProcessor domain.PaymentProcessorServiceService
}

// CreatePaymentSession implements domain.PaymentService.
func (p *paymentService) CreatePaymentSession(ctx context.Context, tripID string, userID string, driverID string, amount int64, currency string) (*types.PaymentIntent, error) {
	metadata := map[string]string{
		"trip_id":   tripID,
		"user_id":   userID,
		"driver_id": driverID,
	}

	sessionID, err := p.paymentProcessor.CreatePaymentSession(ctx, amount, currency, metadata)
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

	return paymentIntent, nil
}

func NewPaymentService(paymentProcessor domain.PaymentProcessorServiceService) domain.PaymentService {
	return &paymentService{
		paymentProcessor: paymentProcessor,
	}
}
