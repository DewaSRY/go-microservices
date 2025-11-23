package models

import (
	"time"

	"gorm.io/gorm"
)

type TripModel struct {
	Id            string         `gorm:"type:varchar(64);primaryKey"`
	Status        string         `gorm:"type:varchar(50);not null;default:'pending'"`
	RiderId       string         `gorm:"type:varchar(64);not null;index"`
	DriverId      string         `gorm:"type:varchar(64);index"`
	TransactionId string         `gorm:"type:varchar(64);index"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

// TableName allows you to customize the table name
func (TripModel) TableName() string {
	return "trips"
}
