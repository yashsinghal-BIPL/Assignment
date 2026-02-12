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
func (s *AccountService) GetAccountByNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	if err := database.DB.Preload("CustomerAccounts").Preload("CustomerAccounts.Customer").Preload("Branch").Preload("Branch.Bank").Where("account_number = ?", accountNumber).First(&account).Error; err != nil {
		return nil, errors.New("account not found")
	}
	return &account, nil
}

// GetCustomerAccounts retrieves all accounts for a customer
func (s *AccountService) GetCustomerAccounts(customerID uint) ([]models.Account, error) {
	var accounts []models.Account
	if err := database.DB.Distinct("accounts.*").
		Joins("JOIN customer_accounts ON customer_accounts.account_id = accounts.id").
		Where("customer_accounts.customer_id = ?", customerID).
		Preload("CustomerAccounts").
		Preload("CustomerAccounts.Customer").
		Preload("Branch").
		Preload("Branch.Bank").
		Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}
	return accounts, nil
}

// Deposit adds money to an account
func (s *AccountService) Deposit(accountID uint, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("deposit amount must be greater than zero")
	}

	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}

	if account.Status != "active" {
		return nil, errors.New("account is not active")
	}

	// Perform transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update account balance
	newBalance := account.Balance + amount
	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:    accountID,
		Type:         models.TransactionTypeDeposit,
		Amount:       amount,
		Description:  description,
		BalanceAfter: newBalance,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}

// Withdraw removes money from an account
func (s *AccountService) Withdraw(accountID uint, amount float64, description string) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("withdrawal amount must be greater than zero")
	}

	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}

	if account.Status != "active" {
		return nil, errors.New("account is not active")
	}

	if account.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	// Perform transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update account balance
	newBalance := account.Balance - amount
	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:    accountID,
		Type:         models.TransactionTypeWithdraw,
		Amount:       amount,
		Description:  description,
		BalanceAfter: newBalance,
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transaction, nil
}
