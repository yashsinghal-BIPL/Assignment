package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	Name             string            `json:"name" gorm:"not null"`
	Email            string            `json:"email" gorm:"uniqueIndex;not null"`
	Phone            string            `json:"phone" gorm:"not null"`
	Address          string            `json:"address"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `json:"-" gorm:"index"`
	CustomerAccounts []CustomerAccount `json:"customer_accounts,omitempty" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	Loans            []Loan            `json:"loans,omitempty" gorm:"foreignKey:CustomerID"`
}
