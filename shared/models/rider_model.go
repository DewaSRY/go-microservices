package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RiderModel struct {
	Id          string         `gorm:"primaryKey;type:varchar(64)"`
	PackageSlug string         `gorm:"type:varchar(24);not null;"`
	IsActive    bool           `gorm:"not null;default:false"`
	Location    datatypes.JSON `gorm:"type:jsonb"`
	Destination datatypes.JSON `gorm:"type:jsonb"`
	UserID      string         `gorm:"type:varchar(64);not null;index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (RiderModel) TableName() string {
	return "riders"
}
