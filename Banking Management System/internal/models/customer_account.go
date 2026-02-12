package models

import (
	"time"

	"gorm.io/gorm"
)

type CustomerAccountRole string

const (
	RolePrimary   CustomerAccountRole = "primary"
	RoleSecondary CustomerAccountRole = "secondary"
	RoleJoint     CustomerAccountRole = "joint"
)

type CustomerAccount struct {
	ID         uint                `json:"id" gorm:"primaryKey"`
	CustomerID uint                `json:"customer_id" gorm:"not null;uniqueIndex:idx_customer_account"`
	Customer   Customer            `json:"customer,omitempty" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE"`
	AccountID  uint                `json:"account_id" gorm:"not null;uniqueIndex:idx_customer_account"`
	Account    Account             `json:"account,omitempty" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	Role       CustomerAccountRole `json:"role" gorm:"type:varchar(20);default:'secondary';not null"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	DeletedAt  gorm.DeletedAt      `json:"-" gorm:"index"`
}

func (CustomerAccount) TableName() string {
	return "customer_accounts"
}
