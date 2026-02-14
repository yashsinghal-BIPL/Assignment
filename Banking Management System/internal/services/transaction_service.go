package services

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"errors"
	"fmt"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

// CreateTransaction creates a deposit or withdraw transaction for an account
func (s *TransactionService) CreateTransaction(accountID uint, amount float64, description string, txType models.TransactionType) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if txType != models.TransactionTypeDeposit && txType != models.TransactionTypeWithdraw {
		return nil, errors.New("type must be deposit or withdraw")
	}

	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}
	if account.Status != "active" {
		return nil, errors.New("account is not active")
	}

	var newBalance float64
	switch txType {
	case models.TransactionTypeDeposit:
		newBalance = account.Balance + amount
	case models.TransactionTypeWithdraw:
		if account.Balance < amount {
			return nil, errors.New("insufficient balance")
		}
		newBalance = account.Balance - amount
	default:
		return nil, errors.New("invalid transaction type")
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	transaction := &models.Transaction{
		AccountID:    accountID,
		Type:         txType,
		Amount:       amount,
		Description:  description,
		BalanceAfter: newBalance,
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit: %w", err)
	}
	return transaction, nil
}

// GetAccountTransactions retrieves all transactions for an account
func (s *TransactionService) GetAccountTransactions(accountID uint, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := database.DB.Where("account_id = ?", accountID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Preload("Loan").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return transactions, nil
}

/* func (s *TransactionService) GetCustomerTransactions(customerID uint, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	// Get all account IDs for the customer
	var accountIDs []uint
	if err := database.DB.Model(&models.Account{}).Where("customer_id = ?", customerID).Pluck("id", &accountIDs).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch account IDs: %w", err)
	}

	if len(accountIDs) == 0 {
		return []models.Transaction{}, nil
	}

	query := database.DB.Where("account_id IN ?", accountIDs).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Preload("Account").Preload("Loan").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return transactions, nil
}
*/
// getting transactions for a loan
func (s *TransactionService) GetLoanTransactions(loanID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Where("loan_id = ?", loanID).Order("created_at DESC").Preload("Account").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch loan transactions: %w", err)
	}
	return transactions, nil
}

// getting transactions by date range for an account
/* func (s *TransactionService) GetTransactionsByDateRange(accountID uint, startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := database.DB.Where("account_id = ? AND created_at BETWEEN ? AND ?", accountID, startDate, endDate).
		Order("created_at DESC").
		Preload("Loan").
		Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}
	return transactions, nil
} */
