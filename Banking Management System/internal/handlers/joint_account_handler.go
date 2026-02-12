package handlers

import (
	"bank-management-system/internal/models"
	"bank-management-system/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JointAccountHandler struct {
	jointAccountService *services.JointAccountService
}

func NewJointAccountHandler(jointAccountService *services.JointAccountService) *JointAccountHandler {
	return &JointAccountHandler{
		jointAccountService: jointAccountService,
	}
}

type CreateJointAccountRequest struct {
	BranchID    uint   `json:"branch_id" binding:"required"`
	CustomerIDs []uint `json:"customer_ids" binding:"required,min=2"`
}

type AddCustomerToAccountRequest struct {
	CustomerID uint   `json:"customer_id" binding:"required"`
	Role       string `json:"role" binding:"required,oneof=primary secondary joint"`
}

type UpdateCustomerRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=primary secondary joint"`
}

// CreateJointAccount creates a new joint account with multiple customers
func (h *JointAccountHandler) CreateJointAccount(c *gin.Context) {
	var req CreateJointAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.jointAccountService.CreateJointAccount(req.BranchID, req.CustomerIDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Joint account created successfully", "account": account})
}

// AddCustomerToAccount adds a customer to an existing account
func (h *JointAccountHandler) AddCustomerToAccount(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	var req AddCustomerToAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.CustomerAccountRole(req.Role)
	link, err := h.jointAccountService.AddCustomerToAccount(req.CustomerID, uint(accountID), role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Customer added to account successfully", "link": link})
}

// RemoveCustomerFromAccount removes a customer from an account
func (h *JointAccountHandler) RemoveCustomerFromAccount(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	err = h.jointAccountService.RemoveCustomerFromAccount(uint(customerID), uint(accountID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer removed from account successfully"})
}

// UpdateCustomerRole updates the role of a customer in an account
func (h *JointAccountHandler) UpdateCustomerRole(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	var req UpdateCustomerRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role := models.CustomerAccountRole(req.Role)
	link, err := h.jointAccountService.UpdateCustomerRole(uint(customerID), uint(accountID), role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer role updated successfully", "link": link})
}

// GetAccountOwners retrieves all owners of an account
func (h *JointAccountHandler) GetAccountOwners(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	owners, err := h.jointAccountService.GetAccountOwners(uint(accountID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, owners)
}

// GetCustomerJointAccounts retrieves all joint accounts for a customer
func (h *JointAccountHandler) GetCustomerJointAccounts(c *gin.Context) {
	customerID, err := strconv.ParseUint(c.Param("customer_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	accounts, err := h.jointAccountService.GetCustomerJointAccounts(uint(customerID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// GetAccountDetails retrieves account details with all owners
func (h *JointAccountHandler) GetAccountDetails(c *gin.Context) {
	accountID, err := strconv.ParseUint(c.Param("account_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	details, err := h.jointAccountService.GetAccountDetails(uint(accountID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, details)
}
