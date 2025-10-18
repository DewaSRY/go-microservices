package service

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/services/payment-service/pkg/types"
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type PaymentStripeProcessService struct {
	config *types.PaymentConfig
}

// CreatePaymentSession implements domain.PaymentProcessorServiceService.
func (p *PaymentStripeProcessService) CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		SuccessURL: stripe.String(p.config.SuccessURL),
		CancelURL:  stripe.String(p.config.CancelURL),
		Metadata:   metadata,
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(currency),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Ride Payment"),
					},
					UnitAmount: stripe.Int64(amount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
	}

	result, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create a payment session on stripe: %w", err)
	}

	return result.ID, nil
}

func NewPaymentStripeProcessService(config *types.PaymentConfig) domain.PaymentProcessorServiceService {
	stripe.Key = config.StripeSecretKey
	return &PaymentStripeProcessService{
		config: config,
	}
}
