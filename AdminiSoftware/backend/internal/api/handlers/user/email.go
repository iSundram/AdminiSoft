
package user

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailHandler struct {
	db           *gorm.DB
	emailService *services.EmailService
}

func NewEmailHandler(db *gorm.DB, emailService *services.EmailService) *EmailHandler {
	return &EmailHandler{
		db:           db,
		emailService: emailService,
	}
}

func (h *EmailHandler) GetEmailAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")
	domainID := c.Param("domain_id")

	var emails []models.Email
	if err := h.db.Where("user_id = ? AND domain_id = ?", userID, domainID).Find(&emails).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}

func (h *EmailHandler) CreateEmailAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		DomainID uint   `json:"domain_id" binding:"required"`
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Quota    int    `json:"quota"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := models.Email{
		UserID:   userID,
		DomainID: req.DomainID,
		Username: req.Username,
		Password: req.Password,
		Quota:    req.Quota,
		Status:   "active",
	}

	if err := h.db.Create(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email account created successfully", "email": email})
}

func (h *EmailHandler) UpdateEmailAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	emailID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var email models.Email
	if err := h.db.Where("id = ? AND user_id = ?", emailID, userID).First(&email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email account not found"})
		return
	}

	var req struct {
		Password string `json:"password"`
		Quota    int    `json:"quota"`
		Status   string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != "" {
		email.Password = req.Password
	}
	if req.Quota > 0 {
		email.Quota = req.Quota
	}
	if req.Status != "" {
		email.Status = req.Status
	}

	if err := h.db.Save(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email account updated successfully", "email": email})
}

func (h *EmailHandler) DeleteEmailAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	emailID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var email models.Email
	if err := h.db.Where("id = ? AND user_id = ?", emailID, userID).First(&email).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email account not found"})
		return
	}

	if err := h.db.Delete(&email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email account deleted successfully"})
}

func (h *EmailHandler) GetEmailForwarders(c *gin.Context) {
	userID := c.GetUint("user_id")
	domainID := c.Param("domain_id")

	var forwarders []models.EmailForwarder
	if err := h.db.Where("user_id = ? AND domain_id = ?", userID, domainID).Find(&forwarders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email forwarders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"forwarders": forwarders})
}

func (h *EmailHandler) CreateEmailForwarder(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		Source      string `json:"source" binding:"required"`
		Destination string `json:"destination" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	forwarder := models.EmailForwarder{
		UserID:      userID,
		DomainID:    req.DomainID,
		Source:      req.Source,
		Destination: req.Destination,
		Status:      "active",
	}

	if err := h.db.Create(&forwarder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email forwarder"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email forwarder created successfully", "forwarder": forwarder})
}
