
package reseller

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountHandler struct {
	db *gorm.DB
}

func NewAccountHandler(db *gorm.DB) *AccountHandler {
	return &AccountHandler{db: db}
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	var accounts []models.User
	if err := h.db.Where("reseller_id = ?", resellerID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}
	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ResellerID = &resellerID
	user.Role = "user"
	user.Status = "active"

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var user models.User
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AccountHandler) SuspendAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var user models.User
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	user.Status = "suspended"
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account suspended successfully"})
}

func (h *AccountHandler) UnsuspendAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var user models.User
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	user.Status = "active"
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsuspend account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account unsuspended successfully"})
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *AccountHandler) GetAccountStats(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	var totalAccounts int64
	var activeAccounts int64
	var suspendedAccounts int64

	h.db.Model(&models.User{}).Where("reseller_id = ?", resellerID).Count(&totalAccounts)
	h.db.Model(&models.User{}).Where("reseller_id = ? AND status = ?", resellerID, "active").Count(&activeAccounts)
	h.db.Model(&models.User{}).Where("reseller_id = ? AND status = ?", resellerID, "suspended").Count(&suspendedAccounts)

	stats := map[string]interface{}{
		"total_accounts":     totalAccounts,
		"active_accounts":    activeAccounts,
		"suspended_accounts": suspendedAccounts,
		"new_accounts_today": 3,
		"new_accounts_week":  15,
		"new_accounts_month": 67,
		"disk_usage":         "145.6 GB",
		"bandwidth_usage":    "2.3 TB",
		"total_domains":      234,
		"total_databases":    156,
		"total_emails":       789,
	}

	c.JSON(http.StatusOK, stats)
}
package reseller

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResellerAccountHandler struct {
	db *gorm.DB
}

func NewResellerAccountHandler(db *gorm.DB) *ResellerAccountHandler {
	return &ResellerAccountHandler{db: db}
}

func (h *ResellerAccountHandler) CreateAccount(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	
	var account models.User
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.ResellerID = &resellerID
	account.Role = "user"
	account.Status = "active"

	if err := h.db.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *ResellerAccountHandler) GetAccounts(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	
	var accounts []models.User
	if err := h.db.Where("reseller_id = ?", resellerID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

func (h *ResellerAccountHandler) UpdateAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resellerID := c.GetUint("user_id")
	
	var account models.User
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *ResellerAccountHandler) SuspendAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resellerID := c.GetUint("user_id")
	
	var account models.User
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	account.Status = "suspended"
	if err := h.db.Save(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to suspend account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account suspended successfully"})
}

func (h *ResellerAccountHandler) DeleteAccount(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	resellerID := c.GetUint("user_id")
	
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).Delete(&models.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}
