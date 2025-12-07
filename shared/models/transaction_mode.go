package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionModel struct {
	Id        string         `gorm:"type:varchar(64);primaryKey"`
	RiderId   string         `gorm:"type:varchar(64);not null;index"`
	DriverId  string         `gorm:"type:varchar(64);not null;index"`
	Status    string         `gorm:"type:varchar(50);not null;default:'pending'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (TransactionModel) TableName() string {
	return "transactions"
}
