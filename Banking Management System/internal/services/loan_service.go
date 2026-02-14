package services

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"errors"
	"fmt"
	"time"
)

const DefaultInterestRate = 12.0

type LoanService struct{}

func NewLoanService() *LoanService {
	return &LoanService{}
}

// CreateLoan creates a new loan for a customer
func (s *LoanService) CreateLoan(customerID, branchID, accountID uint, principalAmount float64) (*models.Loan, error) {
	if principalAmount <= 0 {
		return nil, errors.New("loan amount must be greater than zero")
	}

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

	// Verify account exists and belongs to customer
	var account models.Account
	if err := database.DB.Where("id = ? AND customer_id = ?", accountID, customerID).First(&account).Error; err != nil {
		return nil, errors.New("account not found or does not belong to customer")
	}

	// Generate loan number
	loanNumber := fmt.Sprintf("LOAN%d%d", customerID, time.Now().Unix())

	loan := &models.Loan{
		CustomerID:      customerID,
		BranchID:        branchID,
		AccountID:       accountID,
		LoanNumber:      loanNumber,
		PrincipalAmount: principalAmount,
		RemainingAmount: principalAmount,
		InterestRate:    DefaultInterestRate,
		Status:          models.LoanStatusActive,
		StartDate:       time.Now(),
	}

	if err := database.DB.Create(loan).Error; err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	// Create transaction record for loan disbursement
	transaction := &models.Transaction{
		AccountID:    accountID,
		LoanID:       &loan.ID,
		Type:         models.TransactionTypeLoanTake,
		Amount:       principalAmount,
		Description:  fmt.Sprintf("Loan disbursement - %s", loanNumber),
		BalanceAfter: account.Balance + principalAmount,
	}

	// Update account balance
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&account).Update("balance", account.Balance+principalAmount).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return loan, nil
}

// GetLoan retrieves loan details by loan ID
func (s *LoanService) GetLoan(loanID uint) (*models.Loan, error) {
	_ = s.AccrueInterest(loanID)
	var loan models.Loan
	if err := database.DB.Preload("Customer").Preload("Branch").Preload("Branch.Bank").Preload("Account").First(&loan, loanID).Error; err != nil {
		return nil, errors.New("loan not found")
	}
	return &loan, nil
}

// GetLoanByNumber retrieves loan by loan number
/* func (s *LoanService) GetLoanByNumber(loanNumber string) (*models.Loan, error) {
	var loan models.Loan
	if err := database.DB.Select("id").Where("loan_number = ?", loanNumber).First(&loan).Error; err != nil {
		return nil, errors.New("loan not found")
	}
	_ = s.AccrueInterest(loan.ID)
	if err := database.DB.Preload("Customer").Preload("Branch").Preload("Branch.Bank").Preload("Account").Where("loan_number = ?", loanNumber).First(&loan).Error; err != nil {
		return nil, errors.New("loan not found")
	}
	return &loan, nil
} */

/* // GetCustomerLoans retrieves all loans for a customer
func (s *LoanService) GetCustomerLoans(customerID uint) ([]models.Loan, error) {
	var loans []models.Loan
	if err := database.DB.Preload("Branch").Preload("Branch.Bank").Preload("Account").Where("customer_id = ?", customerID).Find(&loans).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch loans: %w", err)
	}
	return loans, nil
} */

// RepayLoan processes a loan repayment
func (s *LoanService) RepayLoan(loanID, accountID uint, amount float64) (*models.Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("repayment amount must be greater than zero")
	}

	_ = s.AccrueInterest(loanID)

	var loan models.Loan
	if err := database.DB.First(&loan, loanID).Error; err != nil {
		return nil, errors.New("loan not found")
	}

	if loan.Status != models.LoanStatusActive {
		return nil, errors.New("loan is not active")
	}

	var account models.Account
	if err := database.DB.First(&account, accountID).Error; err != nil {
		return nil, errors.New("account not found")
	}

	if account.Balance < amount {
		return nil, errors.New("insufficient balance for repayment")
	}

	// Perform transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate new remaining amount
	newRemainingAmount := loan.RemainingAmount - amount
	if newRemainingAmount < 0 {
		newRemainingAmount = 0
	}

	// Update loan
	updateData := map[string]interface{}{
		"remaining_amount": newRemainingAmount,
	}
	if newRemainingAmount == 0 {
		updateData["status"] = models.LoanStatusPaid
	}

	if err := tx.Model(&loan).Updates(updateData).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update loan: %w", err)
	}

	// Update account balance
	newBalance := account.Balance - amount
	if err := tx.Model(&account).Update("balance", newBalance).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update account balance: %w", err)
	}

	// Create transaction record
	transaction := &models.Transaction{
		AccountID:    accountID,
		LoanID:       &loanID,
		Type:         models.TransactionTypeLoanRepay,
		Amount:       amount,
		Description:  fmt.Sprintf("Loan repayment - %s", loan.LoanNumber),
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

// AccrueInterest updates the loan's remaining amount by adding accrued interest based on elapsed time since last accrual (or StartDate)
func (s *LoanService) AccrueInterest(loanID uint) error {
	var loan models.Loan
	if err := database.DB.First(&loan, loanID).Error; err != nil {
		return errors.New("loan not found")
	}
	if loan.Status != models.LoanStatusActive {
		return nil // No accrual for non-active loans
	}

	// Use UpdatedAt as last accrual date, or StartDate if never accrued
	lastAccrual := loan.UpdatedAt
	if lastAccrual.Before(loan.StartDate) {
		lastAccrual = loan.StartDate
	}
	now := time.Now()
	if !now.After(lastAccrual) {
		return nil // No time has passed
	}

	// Calculate months and years between lastAccrual and now
	years := now.Year() - lastAccrual.Year()
	months := int(now.Month()) - int(lastAccrual.Month())
	days := now.Day() - lastAccrual.Day()
	if days < 0 {
		months--
	}
	if months < 0 {
		years--
		months += 12
	}
	if years < 0 || (years == 0 && months <= 0) {
		return nil // No accrual needed
	}

	// Calculate interest for this period
	principal := loan.RemainingAmount
	rate := loan.InterestRate
	accrued := 0.0
	if years > 0 {
		accrued += principal * rate * float64(years) / 100.0
	}
	if months > 0 {
		accrued += principal * rate * float64(months) / 12.0 / 100.0
	}

	if accrued > 0 {
		if err := database.DB.Model(&loan).Updates(map[string]interface{}{
			"remaining_amount": loan.RemainingAmount + accrued,
			"updated_at":       now,
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

// CalculateInterestForYear calculates the interest to be paid for the current year
func (s *LoanService) CalculateInterestForYear(loanID uint) (float64, error) {
	var loan models.Loan
	if err := database.DB.First(&loan, loanID).Error; err != nil {
		return 0, errors.New("loan not found")
	}

	if loan.Status != models.LoanStatusActive {
		return 0, nil
	}

	// Calculate interest for the year based on remaining amount
	// Interest = (RemainingAmount * InterestRate) / 100
	interestForYear := (loan.RemainingAmount * loan.InterestRate) / 100

	return interestForYear, nil
}

// GetLoanDetails returns comprehensive loan details including pending amount and interest
func (s *LoanService) GetLoanDetails(loanID uint) (map[string]interface{}, error) {
	_ = s.AccrueInterest(loanID)
	loan, err := s.GetLoan(loanID)
	if err != nil {
		return nil, err
	}

	interestForYear, err := s.CalculateInterestForYear(loanID)
	if err != nil {
		return nil, err
	}

	details := map[string]interface{}{
		"loan":              loan,
		"remaining_amount":  loan.RemainingAmount,
		"interest_rate":     loan.InterestRate,
		"interest_for_year": interestForYear,
		"principal_amount":  loan.PrincipalAmount,
		"amount_paid":       loan.PrincipalAmount - loan.RemainingAmount,
		"status":            loan.Status,
		"start_date":        loan.StartDate,
	}

	return details, nil
}
