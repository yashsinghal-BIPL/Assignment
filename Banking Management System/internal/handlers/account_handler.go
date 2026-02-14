package handlers

import (
	"bank-management-system/internal/models"
	"bank-management-system/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService *services.AccountService
}

func NewAccountHandler(accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

type CreateAccountRequest struct {
	BranchID    uint   `json:"branch_id" binding:"required"`
	CustomerIDs []uint `json:"customer_ids" binding:"required,min=1"`
}

type AccountAddCustomerRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	Role       string `json:"role" binding:"required,oneof=primary secondary joint"`
}

// CreateAccount creates a new savings account
func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var (
		account interface{}
		err     error
	)

	// If only one customer is provided, create a regular account.
	// If multiple customers are provided, create a joint account.
	if len(req.CustomerIDs) == 1 {
		account, err = h.accountService.CreateAccount(req.CustomerIDs[0], req.BranchID)
	} else {
		account, err = h.accountService.CreateJointAccount(req.BranchID, req.CustomerIDs)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully", "account": account})
}

// GetAccount retrieves account details
func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	account, err := h.accountService.GetAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// GetCustomerAccounts retrieves all accounts for a customer
/* func (h *AccountHandler) GetCustomerAccounts(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	accounts, err := h.accountService.GetCustomerAccounts(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
} */

// AddCustomerToAccount adds a customer to an existing account
func (h *AccountHandler) AddCustomerToAccount(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var req AccountAddCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.CustomerAccountRole(req.Role)
	link, err := h.accountService.AddCustomerToAccount(req.CustomerID, uint(accountID), role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer added to account successfully", "link": link})
}

// RemoveCustomerFromAccount removes a customer from an account
func (h *AccountHandler) RemoveCustomerFromAccount(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	if err := h.accountService.RemoveCustomerFromAccount(uint(customerID), uint(accountID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer removed from account successfully"})
}
