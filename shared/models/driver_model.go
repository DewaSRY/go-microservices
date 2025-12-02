package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type DriverModel struct {
	Id          string         `gorm:"primaryKey;type:varchar(64)"`
	PackageSlug string         `gorm:"type:varchar(24);not null;"`
	IsActive    bool           `gorm:"not null;default:false"`
	UserId      string         `gorm:"type:varchar(64);not null;index"`
	Location    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (DriverModel) TableName() string {
	return "drivers"
}
