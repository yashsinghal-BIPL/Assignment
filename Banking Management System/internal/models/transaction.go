package models

import (
	"gorm.io/gorm"
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit   TransactionType = "deposit"
	TransactionTypeWithdraw  TransactionType = "withdraw"
	TransactionTypeLoanRepay TransactionType = "loan_repay"
	TransactionTypeLoanTake  TransactionType = "loan_take"
)

type Transaction struct {
	ID           uint            `json:"id" gorm:"primaryKey"`
	AccountID    uint            `json:"account_id" gorm:"not null"`
	Account      Account         `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	LoanID       *uint           `json:"loan_id,omitempty"`
	Loan         *Loan           `json:"loan,omitempty" gorm:"foreignKey:LoanID"`
	Type         TransactionType `json:"type" gorm:"type:varchar(20);not null"`
	Amount       float64         `json:"amount" gorm:"not null"`
	Description  string          `json:"description"`
	BalanceAfter float64         `json:"balance_after"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"-" gorm:"index"`
}
