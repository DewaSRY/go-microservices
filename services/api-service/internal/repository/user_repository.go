package repository

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *db.PostgresManager
}

// CleanUpDriverData implements domain.UserRepository.
func (t *userRepository) CleanUpDriverData(ctx context.Context, connection string) error {

	result := t.db.DB.WithContext(ctx).Model(models.DriverModel{}).
		Where("id = ?", connection).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CleanUpRiderData implements domain.UserRepository.
func (t *userRepository) CleanUpRiderData(ctx context.Context, connection string) error {

	result := t.db.DB.WithContext(ctx).Model(models.RiderModel{}).
		Where("id = ?", connection).
		Update("is_active", false)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GetWaitingRiderIdConnectionList implements domain.UserRepository.
func (t *userRepository) GetWaitingRiderIdConnectionList(ctx context.Context) ([]string, error) {
	var ridersId []string

	result := t.db.DB.WithContext(ctx).
		Model(&models.TripModel{}).
		Select("rider_id").
		Where("LENGTH(transaction_id)  = 0").
		Pluck("rider_id", &ridersId)

	if result.Error != nil {
		return nil, result.Error
	}

	return ridersId, nil

}

// GetDriverList implements domain.UserRepository.
func (t *userRepository) GetDriverActiveList(ctx context.Context) ([]models.DriverModel, error) {
	var drivers []models.DriverModel

	result := t.db.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&drivers)

	if result.Error != nil {
		return nil, result.Error
	}

	return drivers, nil
}

// CreateOrUpdateDriverModel implements domain.UserRepository.
func (t *userRepository) CreateOrUpdateDriverModel(ctx context.Context, model models.DriverModel) error {
	db := t.db.DB.WithContext(ctx)
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"location", "is_active", "updated_at", "package_slug"}),
	}).Create(&model).Error
}

// CreateOrUpdateRiderModel implements domain.UserRepository.
func (t *userRepository) CreateOrUpdateRiderModel(ctx context.Context, model models.RiderModel) error {
	db := t.db.DB.WithContext(ctx)
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"location", "destination", "updated_at"}),
	}).Create(&model).Error
}

// UpdateRiderLocation implements domain.UserRepository.
func (t *userRepository) UpdateRiderLocation(ctx context.Context, riderId string, location []byte, destination []byte) error {
	if result := t.db.DB.WithContext(ctx).Model(&models.RiderModel{}).Where("id = ?", riderId).Updates(
		map[string]interface{}{
			"location":    location,
			"destination": destination,
			"updated_at":  gorm.Expr("NOW()"),
		},
	); result.Error != nil {
		return result.Error
	}
	return nil
}

// CreateDriver implements domain.UserRepository.
func (t *userRepository) CreateDriver(ctx context.Context, connectionId string, data types.CreateDriverParam) error {
	jsonLocation, err := json.Marshal(data.Location)

	if err != nil {
		return err
	}

	newDriver := &models.DriverModel{
		Id:          connectionId,
		PackageSlug: data.PackageSlug,
		IsActive:    true,
		Location:    jsonLocation,
	}

	if result := t.db.DB.WithContext(ctx).Create(newDriver); result.Error != nil {
		return result.Error
	}

	return nil
}

// CreateRider implements domain.UserRepository.
func (t *userRepository) CreateRider(ctx context.Context, connectionId string, data types.CreateRiderParam) error {

	jsonLocation, err := json.Marshal(data.Location)

	if err != nil {
		return err
	}

	newCreateRider := &models.RiderModel{
		Id:       connectionId,
		IsActive: true,
		Location: jsonLocation,
	}

	if result := t.db.DB.WithContext(ctx).Create(newCreateRider); result.Error != nil {
		return result.Error
	}

	return nil
}

func NewUserRepository(db *db.PostgresManager) domain.UserRepository {
	return &userRepository{db: db}
}
