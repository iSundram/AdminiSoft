
package admin

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/services"
	"AdminiSoftware/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AccountHandler struct {
	db             *gorm.DB
	logger         *utils.Logger
	accountService *services.AccountService
}

func NewAccountHandler(db *gorm.DB, logger *utils.Logger, accountService *services.AccountService) *AccountHandler {
	return &AccountHandler{
		db:             db,
		logger:         logger,
		accountService: accountService,
	}
}

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	status := c.Query("status")

	accounts, total, err := h.accountService.ListAccounts(page, limit, search, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	type CreateAccountRequest struct {
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		Domain      string `json:"domain" binding:"required"`
		PackageID   uint   `json:"package_id" binding:"required"`
		ResellerID  uint   `json:"reseller_id"`
		ContactInfo struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Phone     string `json:"phone"`
			Address   string `json:"address"`
		} `json:"contact_info"`
	}

	var req CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.CreateAccount(req.Username, req.Email, req.Password, req.Domain, req.PackageID, req.ResellerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Account created: " + account.Username)
	c.JSON(http.StatusCreated, gin.H{"account": account})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	account, err := h.accountService.GetAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	type UpdateAccountRequest struct {
		Email     string `json:"email"`
		PackageID uint   `json:"package_id"`
		Status    string `json:"status"`
	}

	var req UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account, err := h.accountService.UpdateAccount(uint(id), req.Email, req.PackageID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Account updated: " + account.Username)
	c.JSON(http.StatusOK, gin.H{"account": account})
}

func (h *AccountHandler) SuspendAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	type SuspendRequest struct {
		Reason string `json:"reason" binding:"required"`
	}

	var req SuspendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.accountService.SuspendAccount(uint(id), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Account suspended: " + strconv.FormatUint(id, 10))
	c.JSON(http.StatusOK, gin.H{"message": "Account suspended successfully"})
}

func (h *AccountHandler) UnsuspendAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = h.accountService.UnsuspendAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Account unsuspended: " + strconv.FormatUint(id, 10))
	c.JSON(http.StatusOK, gin.H{"message": "Account unsuspended successfully"})
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	err = h.accountService.DeleteAccount(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Account deleted: " + strconv.FormatUint(id, 10))
	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *AccountHandler) GetAccountStats(c *gin.Context) {
	stats, err := h.accountService.GetAccountStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}
package admin

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

func (h *AccountHandler) ListAccounts(c *gin.Context) {
	var accounts []models.User
	query := h.db.Model(&models.User{})
	
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit
	
	// Filters
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("username ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	
	var total int64
	query.Count(&total)
	
	err := query.Offset(offset).Limit(limit).Find(&accounts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"accounts": accounts,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req struct {
		Username    string `json:"username" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required,min=8"`
		Role        string `json:"role" binding:"required,oneof=admin reseller user"`
		PackageID   *uint  `json:"package_id"`
		Status      string `json:"status" binding:"oneof=active suspended"`
		QuotaLimit  *int64 `json:"quota_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	var existingUser models.User
	if err := h.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	user := models.User{
		Username:   req.Username,
		Email:      req.Email,
		Role:       req.Role,
		Status:     req.Status,
		PackageID:  req.PackageID,
		QuotaLimit: req.QuotaLimit,
	}

	// Hash password
	// user.Password = hashPassword(req.Password)

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"user":    user,
	})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var user models.User
	if err := h.db.Preload("Package").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	var req struct {
		Email       string `json:"email"`
		Role        string `json:"role"`
		Status      string `json:"status"`
		PackageID   *uint  `json:"package_id"`
		QuotaLimit  *int64 `json:"quota_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.PackageID != nil {
		user.PackageID = req.PackageID
	}
	if req.QuotaLimit != nil {
		user.QuotaLimit = req.QuotaLimit
	}

	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"user":    user,
	})
}

func (h *AccountHandler) SuspendAccount(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	if err := h.db.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *AccountHandler) GetAccountStats(c *gin.Context) {
	var stats struct {
		TotalUsers      int64 `json:"total_users"`
		ActiveUsers     int64 `json:"active_users"`
		SuspendedUsers  int64 `json:"suspended_users"`
		ResellerUsers   int64 `json:"reseller_users"`
		AdminUsers      int64 `json:"admin_users"`
	}

	h.db.Model(&models.User{}).Count(&stats.TotalUsers)
	h.db.Model(&models.User{}).Where("status = ?", "active").Count(&stats.ActiveUsers)
	h.db.Model(&models.User{}).Where("status = ?", "suspended").Count(&stats.SuspendedUsers)
	h.db.Model(&models.User{}).Where("role = ?", "reseller").Count(&stats.ResellerUsers)
	h.db.Model(&models.User{}).Where("role = ?", "admin").Count(&stats.AdminUsers)

	c.JSON(http.StatusOK, stats)
}
