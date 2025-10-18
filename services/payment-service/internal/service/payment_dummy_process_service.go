package service

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"context"
	"log"
)

type paymentDummyProcessService struct {
	config *types.PaymentConfig
}

// CreatePaymentSession implements domain.PaymentProcessorServiceService.
func (p *paymentDummyProcessService) CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error) {
	log.Print(
		amount, currency, metadata,
	)
	return "id", nil
}

func NewPaymentProcessService(config *types.PaymentConfig) domain.PaymentProcessorServiceService {
	return &paymentDummyProcessService{
		config: config,
	}
}
