package db

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresManager struct {
	DB *gorm.DB
}

func NewPostgresManager(ctx context.Context, postgresURI string) (*PostgresManager, error) {

	db, err := gorm.Open(postgres.Open(postgresURI), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed_to_connect_to_postgres: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed_to_get_generic_db_handle: %v", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed_to_ping_database: %v", err)
	}

	return &PostgresManager{DB: db}, nil
}

func (t *PostgresManager) Close() {
	if t.DB != nil {
		sqlDB, err := t.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}
