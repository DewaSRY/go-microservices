package repository

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type tripFlowRepository struct {
	db *db.PostgresManager
}

// CreateOrUpdateTransactionModel implements domain.TripFlowRepository.
func (t *tripFlowRepository) CreateOrUpdateTransactionModel(ctx context.Context, model models.TransactionModel) error {
	db := t.db.DB.WithContext(ctx)

	var existingTrip models.TransactionModel
	err := db.First(&existingTrip, "rider_id = ?", model.RiderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&model).Error
		}

		return err
	}

	return db.Model(&existingTrip).Updates(model).Error
}

// CreateOrUpdateRiderTrip implements domain.TripFlowRepository.
func (t *tripFlowRepository) CreateOrUpdateRiderTrip(ctx context.Context, model models.TripModel) error {
	db := t.db.DB.WithContext(ctx)

	var existingTrip models.TripModel
	err := db.First(&existingTrip, "rider_id = ?", model.RiderId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return db.Create(&model).Error
		}

		return err
	}

	return db.Model(&existingTrip).Updates(model).Error
}

// CreateTrip implements domain.TripFlowRepository.
func (t *tripFlowRepository) CreateTrip(ctx context.Context, riderId string) error {
	newTripModel := models.TripModel{
		Id:      uuid.New().String(),
		Status:  "pending",
		RiderId: riderId,
	}

	if result := t.db.DB.WithContext(ctx).Create(&newTripModel); result.Error != nil {
		return result.Error
	}

	return nil
}

func NewTripFlowRepository(db *db.PostgresManager) domain.TripFlowRepository {
	return &tripFlowRepository{db: db}
}
