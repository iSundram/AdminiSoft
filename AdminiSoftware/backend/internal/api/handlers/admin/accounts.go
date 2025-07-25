
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
