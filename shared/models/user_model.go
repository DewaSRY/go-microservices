package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	ID        string         `gorm:"type:varchar(64);primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(100)"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ClientId  string         `gorm:"type:varchar(100)"`
}

func (UserModel) TableName() string {
	return "users"
}
