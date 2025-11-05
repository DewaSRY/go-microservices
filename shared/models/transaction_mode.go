package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionModel struct {
	ID       string `gorm:"type:varchar(64);primaryKey"`
	RiderID  string `gorm:"type:varchar(64);not null;index"`
	DriverID string `gorm:"type:varchar(64);not null;index"`

	Status string `gorm:"type:varchar(50);not null;default:'pending'"`
	// think about rest filed letter
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}
