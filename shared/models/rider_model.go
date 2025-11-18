package models

import "gorm.io/datatypes"

type RiderModel struct {
	ID          string         `gorm:"primaryKey;type:varchar(64)"`
	PackageSlug string         `gorm:"type:varchar(24);not null;"`
	IsActive    bool           `gorm:"not null;default:false"`
	Location    datatypes.JSON `gorm:"type:jsonb"`
	UserID      string         `gorm:"type:varchar(64);not null;index"`
}

func (RiderModel) TableName() string {
	return "riders"
}
