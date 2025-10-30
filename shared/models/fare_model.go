package models

import (
	"gorm.io/datatypes"
)

type FareModel struct {
	ID                string         `gorm:"type:varchar(64);primaryKey"`
	UserID            string         `gorm:"type:varchar(64);not null;index"`
	PackageSlug       string         `gorm:"type:varchar(50);not null;"`
	TotalPriceInCents float64        `gorm:"not null;default:0"`
	Routes            datatypes.JSON `gorm:"type:jsonb"`
}

func (FareModel) TableName() string {
	return "fares"
}
