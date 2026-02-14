package models

import (
	"time"

	"gorm.io/gorm"
)

type LoanStatus string

const (
	LoanStatusActive   LoanStatus = "active"
	LoanStatusPaid     LoanStatus = "paid"
	LoanStatusPending  LoanStatus = "pending"
)

type Loan struct {
	ID                uint           `json:"id" gorm:"primaryKey"`
	CustomerID        uint           `json:"customer_id" gorm:"not null"`
	Customer          Customer       `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	BranchID          uint           `json:"branch_id" gorm:"not null"`
	Branch            Branch         `json:"branch,omitempty" gorm:"foreignKey:BranchID"`
	AccountID         uint           `json:"account_id" gorm:"not null"`
	Account           Account        `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	LoanNumber        string         `json:"loan_number" gorm:"uniqueIndex;not null"`
	PrincipalAmount   float64        `json:"principal_amount" gorm:"not null"`
	RemainingAmount   float64        `json:"remaining_amount" gorm:"not null"`
	InterestRate      float64        `json:"interest_rate" gorm:"default:12.0;not null"`
	Status            LoanStatus     `json:"status" gorm:"type:varchar(20);default:'active'"`
	StartDate         time.Time      `json:"start_date" gorm:"not null"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
	Transactions      []Transaction  `json:"transactions,omitempty" gorm:"foreignKey:LoanID"`
}
