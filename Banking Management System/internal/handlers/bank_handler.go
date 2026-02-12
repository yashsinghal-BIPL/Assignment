package handlers

import (
	"bank-management-system/internal/database"
	"bank-management-system/internal/models"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

type BankHandler struct{}

func NewBankHandler() *BankHandler {
	return &BankHandler{}
}

type CreateBankRequest struct {
	Name    string `json:"name" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Address string `json:"address"`
}

// creates a new bank
func (h *BankHandler) CreateBank(c *gin.Context) {
	var req CreateBankRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bank := &models.Bank{
		Name:    req.Name,
		Code:    req.Code,
		Address: req.Address,
	}

	if err := database.DB.Create(bank).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create bank: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bank created successfully", "bank": bank})
}

// GetBank retrieves bank details
func (h *BankHandler) GetBank(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank ID"})
		return
	}

	var bank models.Bank
	if err := database.DB.Preload("Branches").First(&bank, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank not found"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

// GetAllBanks retrieves all banks
func (h *BankHandler) GetAllBanks(c *gin.Context) {
	var banks []models.Bank
	if err := database.DB.Find(&banks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch banks"})
		return
	}

	c.JSON(http.StatusOK, banks)
}
