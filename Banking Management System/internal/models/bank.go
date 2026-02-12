package models

import (
	"time"

	"gorm.io/gorm"
)

type Bank struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Code      string         `json:"code" gorm:"uniqueIndex;not null"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Branches  []Branch       `json:"branches,omitempty" gorm:"foreignKey:BankID"`
}
