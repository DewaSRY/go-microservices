package repository

import (
	"DewaSRY/go-microservices/services/api-service/internal/domain"
	"DewaSRY/go-microservices/services/api-service/pkg/types"
	"DewaSRY/go-microservices/shared/db"
	"context"
)

type userRepository struct {
	db *db.PostgresManager
}

// CreateDriver implements domain.UserRepository.
func (u *userRepository) CreateDriver(ctx context.Context, data types.CreateDriverParam) error {
	panic("unimplemented")
}

// CreateRider implements domain.UserRepository.
func (u *userRepository) CreateRider(ctx context.Context, data types.CreateRiderParam) error {
	panic("unimplemented")
}

func NewUserRepository(db *db.PostgresManager) domain.UserRepository {
	return &userRepository{db: db}
}
