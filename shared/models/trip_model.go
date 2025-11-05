package models

import (
	"time"

	"gorm.io/gorm"
)

type TripModel struct {
	ID         string         `gorm:"type:varchar(64);primaryKey"`
	UserID     string         `gorm:"type:varchar(64);not null;index"`
	Status     string         `gorm:"type:varchar(50);not null;default:'pending'"`
	RideFareID string         `gorm:"type:varchar(64);index"`
	DriverID   string         `gorm:"type:varchar(64);index"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

// TableName allows you to customize the table name
func (TripModel) TableName() string {
	return "trips"
}

// func (t *TripModel) ToTripProto() *tripgrpc.Trip {
// 	return &tripgrpc.Trip{
// 		Id:     fmt.Sprintf("{%d}", t.ID),
// 		UserID: fmt.Sprintf("{%d}", t.UserID),
// 		Status:       t.Status,
// 		Route:        t.RideFare.Route.ToRouteProto(),
// 		SelectedFare: t.RideFare.ToTripProto(),
// 	}
// }
