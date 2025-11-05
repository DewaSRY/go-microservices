package repository

import (
	"DewaSRY/go-microservices/services/driver-service/internal/domain"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"fmt"
)

type driverRepository struct {
	db *db.PostgresManager
}

// GetDriverById implements domain.DriverRepository.
func (t *driverRepository) GetDriverById(ctx context.Context, driverId string) (*models.DriverModel, error) {
	var driver models.DriverModel

	if result := t.db.DB.WithContext(ctx).Where("id = ?", driverId).First(&driver); result.Error != nil {
		return nil, fmt.Errorf("failed_to_get_driver_with_id_%s:%w", driverId, result.Error)
	}

	return &driver, nil
}

// UpdateDriverData implements domain.DriverRepository.
func (t *driverRepository) UpdateDriverData(ctx context.Context, driverId string, partialData map[string]interface{}) error {
	if result := t.db.DB.WithContext(ctx).Model(&models.DriverModel{}).Where("id = ?", driverId).Select("*").Updates(partialData); result.Error != nil {
		return fmt.Errorf("failed_driver_with_id_%s:%w", driverId, result.Error)
	}
	return nil
}

// CreateDriver implements domain.DriverRepository.
func (t *driverRepository) CreateDriver(ctx context.Context, driver *models.DriverModel) error {
	if result := t.db.DB.WithContext(ctx).Create(driver); result.Error != nil {
		return fmt.Errorf("failed_to_create_driver:%w", result.Error)
	}
	return nil
}

// GetActiveDriverIdList implements domain.DriverRepository.
func (t *driverRepository) GetActiveDriverIdList(ctx context.Context, filter domain.DbFilter) ([]string, error) {
	var ids []string

	if result := filter(t.db.DB.WithContext(ctx).Model(&models.DriverModel{})).Pluck("id", &ids); result.Error != nil {
		return nil, fmt.Errorf("failed_to_get_driver_id: %w", result.Error)
	}

	return ids, nil
}

func NewDriverRepository(db *db.PostgresManager) domain.DriverRepository {
	return &driverRepository{db: db}
}
