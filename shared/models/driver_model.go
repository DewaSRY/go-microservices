package models

import "gorm.io/datatypes"

type DriverModel struct {
	ID             string         `gorm:"primaryKey;type:varchar(64)"`
	Name           string         `gorm:"type:varchar(50)"`
	Geohash        string         `gorm:"type:varchar(64)"`
	ProfilePicture string         `gorm:"type:varchar(50)"`
	CarPlate       string         `gorm:"type:varchar(50);not null;"`
	PackageSlug    string         `gorm:"type:varchar(24);not null;"`
	IsActive       bool           `gorm:"not null;default:false"`
	Location       datatypes.JSON `gorm:"type:jsonb"`
}

func (DriverModel) TableName() string {
	return "drivers"
}
