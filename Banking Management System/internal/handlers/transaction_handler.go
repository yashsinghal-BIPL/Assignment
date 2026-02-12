package handlers

import (
	"bank-management-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// GetAccountTransactions retrieves transactions for an account
func (h *TransactionHandler) GetAccountTransactions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	limit := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
	}

	transactions, err := h.transactionService.GetAccountTransactions(uint(id), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetCustomerTransactions retrieves all transactions for a customer
func (h *TransactionHandler) GetCustomerTransactions(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	limit := 0
	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
			return
		}
	}

	transactions, err := h.transactionService.GetCustomerTransactions(uint(customerID), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetLoanTransactions retrieves transactions for a loan
func (h *TransactionHandler) GetLoanTransactions(c *gin.Context) {
	loanID, err := strconv.ParseUint(c.Param("loan_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid loan ID"})
		return
	}

	transactions, err := h.transactionService.GetLoanTransactions(uint(loanID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
