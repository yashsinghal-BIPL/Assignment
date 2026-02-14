package models

import (
	"time"

	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeSavings AccountType = "savings"
)

type AccountStatus string

const (
	AccountStatusActive   AccountStatus = "active"
	AccountStatusInactive AccountStatus = "inactive"
	AccountStatusClosed   AccountStatus = "closed"
)

type Account struct {
	ID               uint              `json:"id" gorm:"primaryKey"`
	CustomerID       uint              `json:"customer_id" gorm:"not null"`
	Customer         Customer          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	BranchID         uint              `json:"branch_id" gorm:"not null"`
	Branch           Branch            `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	AccountNumber    string            `json:"account_number" gorm:"uniqueIndex;not null"`
	Type             AccountType       `json:"type" gorm:"type:varchar(20);not null"`
	Balance          float64           `json:"balance" gorm:"default:0;not null"`
	Status           AccountStatus     `json:"status" gorm:"default:'active'"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	DeletedAt        gorm.DeletedAt    `json:"-" gorm:"index"`
	Transactions     []Transaction     `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
	CustomerAccounts []CustomerAccount `json:"customer_accounts,omitempty" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
}
