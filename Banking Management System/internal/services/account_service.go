package services

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"errors"
	"fmt"
	"time"
)

type AccountService struct{}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) CreateAccount(customerID, branchID uint) (*models.Account, error) {
	// Verify customer exists
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		return nil, errors.New("customer not found")
	}

	// Verify branch exists
	var branch models.Branch
	if err := database.DB.First(&branch, branchID).Error; err != nil {
		return nil, errors.New("branch not found")
	}

	// Generate account number
	accountNumber := fmt.Sprintf("ACC%d%d", customerID, time.Now().Unix())

	account := &models.Account{
		CustomerID:    customerID,
		BranchID:      branchID,
		AccountNumber: accountNumber,
		Type:          models.AccountTypeSavings,
		Balance:       0,
		Status:        "active",
	}

	// Use transaction to ensure consistency
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	// Link customer as primary owner
	link := &models.CustomerAccount{
		CustomerID: customerID,
		AccountID:  account.ID,
		Role:       models.RolePrimary,
	}
	if err := tx.Create(link).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to link customer: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return account, nil
}

// GetAccount retrieves account details by account ID
func (s *AccountService) GetAccount(accountID uint) (*models.Account, error) {
	var account models.Account
	if err := database.DB.Preload("CustomerAccounts").Preload("CustomerAccounts.Customer").Preload("Branch").Preload("Branch.Bank").First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}
	return &account, nil
}

// GetAccountByNumber retrieves account by account number
/* func (s *AccountService) GetAccountByNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	if err := database.DB.Preload("CustomerAccounts").Preload("CustomerAccounts.Customer").Preload("Branch").Preload("Branch.Bank").Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return nil, errors.New("account not found")
	}
	return &account, nil
} */

// CreateJointAccount creates a new joint account with multiple customers
func (s *AccountService) CreateJointAccount(branchID uint, customerIDs []uint) (*models.Account, error) {
	if len(customerIDs) < 2 {
		return nil, errors.New("joint account must have at least 2 customers")
	}

	var customers []models.Customer
	if err := database.DB.Where("id IN ?", customerIDs).Find(&customers).Error; err != nil {
		return nil, errors.New("failed to verify customers")
	}
	if len(customers) != len(customerIDs) {
		return nil, errors.New("one or more customers not found")
	}

	var branch models.Branch
	if err := database.DB.First(&branch, branchID).Error; err != nil {
		return nil, errors.New("branch not found")
	}

	accountNumber := fmt.Sprintf("JOINT%d%d", branchID, time.Now().Unix())
	account := &models.Account{
		CustomerID:    customerIDs[0],
		BranchID:      branchID,
		AccountNumber: accountNumber,
		Type:          models.AccountTypeSavings,
		Balance:       0,
		Status:        "active",
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(account).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	primaryLink := &models.CustomerAccount{
		CustomerID: customerIDs[0],
		AccountID:  account.ID,
		Role:       models.RolePrimary,
	}
	if err := tx.Create(primaryLink).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to link primary customer: %w", err)
	}

	for i := 1; i < len(customerIDs); i++ {
		jointLink := &models.CustomerAccount{
			CustomerID: customerIDs[i],
			AccountID:  account.ID,
			Role:       models.RoleJoint,
		}
		if err := tx.Create(jointLink).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to link customer: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return account, nil
}

// AddCustomerToAccount adds a customer to an existing account
func (s *AccountService) AddCustomerToAccount(customerID, accountID uint, role models.CustomerAccountRole) (*models.CustomerAccount, error) {
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		return nil, errors.New("customer not found")
	}
	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}
	var existingLink models.CustomerAccount
	if err := database.DB.Where("customer_id = ? AND account_id = ?", customerID, accountID).First(&existingLink).Error; err == nil {
		return nil, errors.New("customer is already linked to this account")
	}
	link := &models.CustomerAccount{
		CustomerID: customerID,
		AccountID:  accountID,
		Role:       role,
	}
	if err := database.DB.Create(link).Error; err != nil {
		return nil, fmt.Errorf("failed to add customer to account: %w", err)
	}
	return link, nil
}

// RemoveCustomerFromAccount removes a customer from an account
func (s *AccountService) RemoveCustomerFromAccount(customerID, accountID uint) error {
	var count int64
	if err := database.DB.Model(&models.CustomerAccount{}).
		Where("account_id = ? AND customer_id != ?", accountID, customerID).
		Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check remaining customers: %w", err)
	}
	if count == 0 {
		return errors.New("cannot remove the last customer from an account")
	}
	if err := database.DB.Where("customer_id = ? AND account_id = ?", customerID, accountID).
		Delete(&models.CustomerAccount{}).Error; err != nil {
		return fmt.Errorf("failed to remove customer from account: %w", err)
	}
	return nil
}

// GetAccountOwners returns all owners of an account
func (s *AccountService) GetAccountOwners(accountID uint) ([]models.CustomerAccount, error) {
	var owners []models.CustomerAccount
	if err := database.DB.Preload("Customer").
		Where("account_id = ?", accountID).
		Order("role DESC, created_at ASC").
		Find(&owners).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch account owners: %w", err)
	}
	return owners, nil
}

// AccountWithOwners represents account with its owners
type AccountWithOwners struct {
	Account models.Account           `json:"account"`
	Owners  []models.CustomerAccount `json:"owners"`
}

// GetAccountDetails returns account details with all owners
func (s *AccountService) GetAccountDetails(accountID uint) (*AccountWithOwners, error) {
	var account models.Account
	if err := database.DB.Preload("Branch").Preload("Branch.Bank").First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}
	owners, err := s.GetAccountOwners(accountID)
	if err != nil {
		return nil, err
	}
	return &AccountWithOwners{
		Account: account,
		Owners:  owners,
	}, nil
}
