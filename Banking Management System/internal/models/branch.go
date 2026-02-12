package models

import (
	"time"

	"gorm.io/gorm"
)

type Branch struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	BankID    uint           `json:"bank_id" gorm:"not null"`
	Bank      Bank           `json:"bank,omitempty" gorm:"foreignKey:BankID"`
	Name      string         `json:"name" gorm:"not null"`
	Code      string         `json:"code" gorm:"uniqueIndex;not null"`
	Address   string         `json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Accounts  []Account      `json:"accounts,omitempty" gorm:"foreignKey:BranchID"`
	Loans     []Loan         `json:"loans,omitempty" gorm:"foreignKey:BranchID"`
}
