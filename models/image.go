package models

import (
	"time"

	"gorm.io/gorm"
)

type Image struct {
	ID        uint           `gorm:"primaryKey"`
	FileName  string         `gorm:"not null"`
	URL       string         `gorm:"not null"`
	Title     string         `gorm:"not null"`
	UserID    uint           `gorm:"not null"`
	User      User           `gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
