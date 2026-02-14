package handlers

import (
	"bank-management-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LoanHandler struct {
	loanService *services.LoanService
}

func NewLoanHandler(loanService *services.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

type CreateLoanRequest struct {
	CustomerID      uint    `json:"customer_id" binding:"required"`
	BranchID        uint    `json:"branch_id" binding:"required"`
	AccountID       uint    `json:"account_id" binding:"required"`
	PrincipalAmount float64 `json:"principal_amount" binding:"required,gt=0"`
}

type RepayLoanRequest struct {
	AccountID uint    `json:"account_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
}

// CreateLoan creates a new loan
func (h *LoanHandler) CreateLoan(c *gin.Context) {
	var req CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.loanService.CreateLoan(req.CustomerID, req.BranchID, req.AccountID, req.PrincipalAmount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Loan created successfully", "loan": loan})
}

// GetLoan retrieves loan details
/* func (h *LoanHandler) GetLoan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	loan, err := h.loanService.GetLoan(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loan)
} */

// GetLoanDetails retrieves comprehensive loan details including pending amount and interest
func (h *LoanHandler) GetLoanDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	details, err := h.loanService.GetLoanDetails(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}

// GetCustomerLoans retrieves all loans for a customer
/* func (h *LoanHandler) GetCustomerLoans(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	loans, err := h.loanService.GetCustomerLoans(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, loans)
} */

// RepayLoan handles loan repayment
func (h *LoanHandler) RepayLoan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	var req RepayLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.loanService.RepayLoan(uint(id), req.AccountID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan repayment successful", "transaction": transaction})
}
