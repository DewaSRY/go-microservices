package repository

import (
	"DewaSRY/go-microservices/services/trip-service/internal/domain"
	"DewaSRY/go-microservices/shared/db"
	drivergrpc "DewaSRY/go-microservices/shared/proto/driver_proto"
	"context"
	"fmt"

	"DewaSRY/go-microservices/shared/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type tripRepository struct {
	db *db.PostgresManager
}

// FindTripWithFilter implements domain.TripRepository.
func (t *tripRepository) FindTripWithFilter(dbFilter domain.DbFilter) (*models.TripModel, error) {
	var trip models.TripModel
	if result := dbFilter(t.db.DB).Find(&trip); result.Error != nil {
		return nil, result.Error
	}
	if len(trip.ID) == 0 {
		return nil, fmt.Errorf("trip_not_found")
	}

	return &trip, nil
}

// GetFareById implements domain.TripRepository.
func (t *tripRepository) GetFareById(ctx context.Context, rideFareID string) (*models.FareModel, error) {
	var fare models.FareModel
	result := t.db.DB.WithContext(ctx).Where("id = ?", rideFareID).First(&fare)
	if result.Error != nil {
		return nil, fmt.Errorf("failed_to_get_ride_by_id: %w", result.Error)
	}

	return &fare, nil
}

// GetTripByID implements domain.TripRepository.
func (t *tripRepository) GetTripByID(ctx context.Context, id string) (*models.TripModel, error) {
	var trip models.TripModel
	result := t.db.DB.WithContext(ctx).Where("id = ?", id).First(&trip)
	if result.Error != nil {
		return nil, fmt.Errorf("failed_to_get_trip_by_id: %w", result.Error)
	}

	return &trip, nil
}

// CreateTrip implements domain.TripRepository.
func (t *tripRepository) CreateTrip(ctx context.Context, trip *models.TripModel) error {
	if result := t.db.DB.WithContext(ctx).Create(trip); result.Error != nil {
		return fmt.Errorf("failed_to_create_trip: %w", result.Error)
	}
	return nil
}

// CreateTrip implements domain.TripRepository.
func (t *tripRepository) CreateFare(ctx context.Context, fare *models.FareModel) error {
	if result := t.db.DB.WithContext(ctx).Create(fare); result.Error != nil {
		return fmt.Errorf("failed_to_create_fare: %v", result.Error)
	}
	return nil
}

// CreateFareList implements domain.TripRepository.
func (t *tripRepository) CreateFareList(ctx context.Context, fares []*models.FareModel) error {
	if len(fares) == 0 {
		return nil
	}

	// Use transaction for bulk operations
	err := t.db.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, fare := range fares {
			// Use Clauses with OnConflict for upsert behavior
			result := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"user_id",
					"package_slug",
					"total_price_in_cents",
					"routes",
				}),
			}).Create(fare)

			if result.Error != nil {
				return fmt.Errorf("failed_to_save_fare_list: %w", result.Error)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed_to_save_fare_lis: %w", err)
	}

	return nil
}

// update this letter
// UpdateTrip implements domain.TripRepository.
func (t *tripRepository) UpdateTrip(ctx context.Context, tripID string, status string, driver *drivergrpc.Driver) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if driver != nil {
		updates["driver_id"] = driver.Id
	}

	result := t.db.DB.WithContext(ctx).
		Model(&models.TripModel{}).
		Where("id = ?", tripID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed_to_update_trip: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("trip_not_found_with_id: %s", tripID)
	}

	return nil
}

func NewTripRepository(db *db.PostgresManager) domain.TripRepository {
	return &tripRepository{db: db}
}
