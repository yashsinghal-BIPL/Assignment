package services

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"errors"
	"fmt"
	"time"
)

type JointAccountService struct{}

func NewJointAccountService() *JointAccountService {
	return &JointAccountService{}
}

// CreateJointAccount creates a new joint account with multiple customers
func (s *JointAccountService) CreateJointAccount(branchID uint, customerIDs []uint) (*models.Account, error) {
	if len(customerIDs) < 2 {
		return nil, errors.New("joint account must have at least 2 customers")
	}

	// Verify all customers exist
	var customers []models.Customer
	if err := database.DB.Where("id IN ?", customerIDs).Find(&customers).Error; err != nil {
		return nil, errors.New("failed to verify customers")
	}

	if len(customers) != len(customerIDs) {
		return nil, errors.New("one or more customers not found")
	}

	// Verify branch exists
	var branch models.Branch
	if err := database.DB.First(&branch, branchID).Error; err != nil {
		return nil, errors.New("branch not found")
	}

	// Generate account number
	accountNumber := fmt.Sprintf("JOINT%d%d", branchID, time.Now().Unix())

	// Create account
	account := &models.Account{
		CustomerID:    customerIDs[0],
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

	// Add first customer as primary owner
	primaryLink := &models.CustomerAccount{
		CustomerID: customerIDs[0],
		AccountID:  account.ID,
		Role:       models.RolePrimary,
	}
	if err := tx.Create(primaryLink).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to link primary customer: %w", err)
	}

	// Add remaining customers as joint owners
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
func (s *JointAccountService) AddCustomerToAccount(customerID, accountID uint, role models.CustomerAccountRole) (*models.CustomerAccount, error) {
	// Verify customer exists
	var customer models.Customer
	if err := database.DB.First(&customer, customerID).Error; err != nil {
		return nil, errors.New("customer not found")
	}

	// Verify account exists
	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}

	// Check if link already exists
	var existingLink models.CustomerAccount
	if err := database.DB.Where("customer_id = ? AND account_id = ?", customerID, accountID).
		First(&existingLink).Error; err == nil {
		return nil, errors.New("customer is already linked to this account")
	}

	// Create the link
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
func (s *JointAccountService) RemoveCustomerFromAccount(customerID, accountID uint) error {
	// Check that account will still have at least one customer
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

// UpdateCustomerRole updates the role of a customer in an account
func (s *JointAccountService) UpdateCustomerRole(customerID, accountID uint, newRole models.CustomerAccountRole) (*models.CustomerAccount, error) {
	// Find the link
	var link models.CustomerAccount
	if err := database.DB.Where("customer_id = ? AND account_id = ?", customerID, accountID).
		First(&link).Error; err != nil {
		return nil, errors.New("customer-account link not found")
	}

	// Update the role
	if err := database.DB.Model(&link).Update("role", newRole).Error; err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &link, nil
}

// GetAccountOwners retrieves all owners of an account
func (s *JointAccountService) GetAccountOwners(accountID uint) ([]models.CustomerAccount, error) {
	var owners []models.CustomerAccount
	if err := database.DB.Preload("Customer").
		Where("account_id = ?", accountID).
		Order("role DESC, created_at ASC").
		Find(&owners).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch account owners: %w", err)
	}
	return owners, nil
}

// GetCustomerJointAccounts retrieves all joint accounts for a customer
func (s *JointAccountService) GetCustomerJointAccounts(customerID uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := database.DB.Distinct("accounts.*").
		Joins("JOIN customer_accounts ON customer_accounts.account_id = accounts.id").
		Where("customer_accounts.customer_id = ?", customerID).
		Preload("Branch").
		Preload("Branch.Bank").
		Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}
	return accounts, nil
}

// GetAccountDetails retrieves account details with all owners
func (s *JointAccountService) GetAccountDetails(accountID uint) (*AccountWithOwners, error) {
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

// AccountWithOwners represents account with its owners
type AccountWithOwners struct {
	Account models.Account           `json:"account"`
	Owners  []models.CustomerAccount `json:"owners"`
}
