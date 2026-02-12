package services

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"fmt"
	//"time"
)

type TransactionService struct{}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
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

func (s *TransactionService) GetCustomerTransactions(customerID uint, limit int) ([]models.Transaction, error) {
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
