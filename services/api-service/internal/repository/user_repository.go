package repository

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/db"
	"DewaSRY/go-microservices/shared/models"
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type userRepository struct {
	db *db.PostgresManager
}

// CreateDriver implements domain.UserRepository.
func (t *userRepository) CreateDriver(ctx context.Context, data types.CreateDriverParam) error {
	jsonLocation, err := json.Marshal(data.Location)

	if err != nil {
		return err
	}

	newDriver := models.DriverModel{
		ID:          uuid.New().String(),
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
func (t *userRepository) CreateRider(ctx context.Context, data types.CreateRiderParam) error {

	jsonLocation, err := json.Marshal(data.Location)

	if err != nil {
		return err
	}

	newCreateRider := models.RiderModel{
		ID:          uuid.New().String(),
		PackageSlug: data.PackageSlug,
		IsActive:    true,
		Location:    jsonLocation,
	}

	if result := t.db.DB.WithContext(ctx).Create(newCreateRider); result.Error != nil {
		return result.Error
	}

	return nil
}

func NewUserRepository(db *db.PostgresManager) domain.UserRepository {
	return &userRepository{db: db}
}
