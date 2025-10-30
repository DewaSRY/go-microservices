package models

type DriverModel struct {
	ID             string `gorm:"primaryKey;type:varchar(64)"`
	Name           string `gorm:"type:varchar(50);not null;"`
	ProfilePicture string `gorm:"type:varchar(50);not null;"`
	CarPlate       string `gorm:"type:varchar(50);not null;"`
	ClientId       string `gorm:"type:varchar(100)"`
}

func (DriverModel) TableName() string {
	return "drivers"
}
