package handlers

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BranchHandler struct{}

func NewBranchHandler() *BranchHandler {
	return &BranchHandler{}
}

type CreateBranchRequest struct {
	BankID  uint   `json:"bank_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Address string `json:"address"`
}

// CreateBranch creates a new branch
func (h *BranchHandler) CreateBranch(c *gin.Context) {
	var req CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify bank exists
	var bank models.Bank
	if err := database.DB.First(&bank, req.BankID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bank not found"})
		return
	}

	branch := &models.Branch{
		BankID:  req.BankID,
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
	}

	if err := database.DB.Create(branch).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create branch: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Branch created successfully", "branch": branch})
}

// GetBranch retrieves branch details
func (h *BranchHandler) GetBranch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid branch ID"})
		return
	}

	var branch models.Branch
	if err := database.DB.Preload("Bank").First(&branch, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "branch not found"})
		return
	}

	c.JSON(http.StatusOK, branch)
}

// GetBranchesByBank retrieves all branches for a bank
func (h *BranchHandler) GetBranchesByBank(c *gin.Context) {
	bankID, err := strconv.ParseUint(c.Param("bank_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank ID"})
		return
	}

	var branches []models.Branch
	if err := database.DB.Where("bank_id = ?", bankID).Find(&branches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch branches"})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// GetAllBranches retrieves all branches
func (h *BranchHandler) GetAllBranches(c *gin.Context) {
	var branches []models.Branch
	if err := database.DB.Preload("Bank").Find(&branches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch branches"})
		return
	}

	c.JSON(http.StatusOK, branches)
}
