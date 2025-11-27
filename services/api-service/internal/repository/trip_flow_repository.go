package repository

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

type tripFlowRepository struct {
	db *db.PostgresManager
}

// CreateOrUpdateRiderTrip implements domain.TripFlowRepository.
func (t *tripFlowRepository) CreateOrUpdateRiderTrip(ctx context.Context, model models.TripModel) error {
	db := t.db.DB.WithContext(ctx)
	return db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "rider_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
		},
	).Create(&model).Error
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
