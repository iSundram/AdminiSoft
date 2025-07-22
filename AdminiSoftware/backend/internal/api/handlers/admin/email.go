
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailHandler struct {
	db *gorm.DB
}

func NewEmailHandler(db *gorm.DB) *EmailHandler {
	return &EmailHandler{db: db}
}

func (h *EmailHandler) GetEmailAccounts(c *gin.Context) {
	var accounts []models.EmailAccount
	if err := h.db.Preload("User").Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *EmailHandler) CreateEmailAccount(c *gin.Context) {
	var account models.EmailAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.Status = "active"
	if err := h.db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *EmailHandler) UpdateEmailAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var account models.EmailAccount
	if err := h.db.First(&account, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email account not found"})
		return
	}

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *EmailHandler) DeleteEmailAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.EmailAccount{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete email account"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email account deleted successfully"})
}

func (h *EmailHandler) GetForwarders(c *gin.Context) {
	var forwarders []models.EmailForwarder
	if err := h.db.Preload("User").Find(&forwarders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch email forwarders"})
		return
	}
	c.JSON(http.StatusOK, forwarders)
}

func (h *EmailHandler) CreateForwarder(c *gin.Context) {
	var forwarder models.EmailForwarder
	if err := c.ShouldBindJSON(&forwarder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	forwarder.Status = "active"
	if err := h.db.Create(&forwarder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email forwarder"})
		return
	}

	c.JSON(http.StatusCreated, forwarder)
}

func (h *EmailHandler) GetMailQueue(c *gin.Context) {
	queue := []map[string]interface{}{
		{
			"id":         1,
			"from":       "system@example.com",
			"to":         "user@example.com",
			"subject":    "Welcome to AdminiSoftware",
			"size":       "2.5KB",
			"status":     "queued",
			"attempts":   0,
			"created_at": "2024-01-15T10:30:00Z",
		},
		{
			"id":         2,
			"from":       "noreply@example.com",
			"to":         "admin@example.com",
			"subject":    "System notification",
			"size":       "1.2KB",
			"status":     "failed",
			"attempts":   3,
			"created_at": "2024-01-15T09:45:00Z",
		},
	}
	c.JSON(http.StatusOK, queue)
}

func (h *EmailHandler) GetMailStats(c *gin.Context) {
	stats := map[string]interface{}{
		"sent_today":     245,
		"sent_this_week": 1680,
		"sent_this_month": 6420,
		"queued":         3,
		"failed":         12,
		"bounced":        8,
		"spam_blocked":   67,
		"disk_usage":     "2.4GB",
		"quota":          "10GB",
	}
	c.JSON(http.StatusOK, stats)
}

func (h *EmailHandler) GetSpamSettings(c *gin.Context) {
	settings := map[string]interface{}{
		"spam_assassin_enabled":  true,
		"spam_threshold":         5.0,
		"auto_delete_spam":       false,
		"quarantine_enabled":     true,
		"whitelist":              []string{"trusted@example.com"},
		"blacklist":              []string{"spam@example.com"},
		"greylisting_enabled":    true,
		"rbl_checks_enabled":     true,
		"dkim_enabled":           true,
		"spf_enabled":            true,
		"dmarc_enabled":          true,
	}
	c.JSON(http.StatusOK, settings)
}

func (h *EmailHandler) UpdateSpamSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Spam settings updated successfully"})
}
