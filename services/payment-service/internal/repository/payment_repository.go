package repository

import (
	"DewaSRY/go-microservices/services/payment-service/internal/domain"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"fmt"
)

type paymentRepository struct {
	db *db.PostgresManager
}

// CreateTransaction implements domain.PaymentRepo.
func (t *paymentRepository) CreateTransaction(ctx context.Context, transaction *models.TransactionModel) error {
	if result := t.db.DB.WithContext(ctx).Create(transaction); result.Error != nil {
		return fmt.Errorf("failed_to_create_transaction: %w", result.Error)
	}
	return nil
}

func NewPaymentRepository(db *db.PostgresManager) domain.PaymentRepo {
	return &paymentRepository{db: db}
}
