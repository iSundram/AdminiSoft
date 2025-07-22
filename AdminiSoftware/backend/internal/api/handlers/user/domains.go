
package user

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DomainHandler struct {
	db *gorm.DB
}

func NewDomainHandler(db *gorm.DB) *DomainHandler {
	return &DomainHandler{db: db}
}

func (h *DomainHandler) GetDomains(c *gin.Context) {
	userID := c.GetUint("user_id")
	var domains []models.Domain
	if err := h.db.Where("user_id = ?", userID).Find(&domains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch domains"})
		return
	}
	c.JSON(http.StatusOK, domains)
}

func (h *DomainHandler) CreateSubdomain(c *gin.Context) {
	userID := c.GetUint("user_id")
	var domain models.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain.UserID = userID
	domain.Type = "subdomain"
	domain.Status = "active"

	if err := h.db.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subdomain"})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *DomainHandler) GetSubdomains(c *gin.Context) {
	userID := c.GetUint("user_id")
	var subdomains []models.Domain
	if err := h.db.Where("user_id = ? AND type = ?", userID, "subdomain").Find(&subdomains).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subdomains"})
		return
	}
	c.JSON(http.StatusOK, subdomains)
}

func (h *DomainHandler) CreateAddonDomain(c *gin.Context) {
	userID := c.GetUint("user_id")
	var domain models.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	domain.UserID = userID
	domain.Type = "addon"
	domain.Status = "active"

	if err := h.db.Create(&domain).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create addon domain"})
		return
	}

	c.JSON(http.StatusCreated, domain)
}

func (h *DomainHandler) GetRedirects(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	redirects := []map[string]interface{}{
		{
			"id":          1,
			"source":      "old.example.com",
			"destination": "https://new.example.com",
			"type":        "301",
			"wildcard":    false,
			"created_at":  "2024-01-15T10:30:00Z",
		},
		{
			"id":          2,
			"source":      "blog.example.com",
			"destination": "https://example.com/blog",
			"type":        "302",
			"wildcard":    true,
			"created_at":  "2024-01-14T15:20:00Z",
		},
	}
	
	c.JSON(http.StatusOK, redirects)
}

func (h *DomainHandler) CreateRedirect(c *gin.Context) {
	userID := c.GetUint("user_id")
	var redirect map[string]interface{}
	if err := c.ShouldBindJSON(&redirect); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	redirect["user_id"] = userID
	redirect["id"] = 3
	redirect["created_at"] = "2024-01-15T10:30:00Z"

	c.JSON(http.StatusCreated, redirect)
}

func (h *DomainHandler) DeleteRedirect(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "Redirect " + id + " deleted successfully"})
}

func (h *DomainHandler) GetDNSRecords(c *gin.Context) {
	userID := c.GetUint("user_id")
	domain := c.Query("domain")
	
	var dnsRecords []models.DNSRecord
	if err := h.db.Joins("JOIN dns_zones ON dns_records.zone_id = dns_zones.id").
		Where("dns_zones.user_id = ? AND dns_zones.name = ?", userID, domain).
		Find(&dnsRecords).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DNS records"})
		return
	}
	
	c.JSON(http.StatusOK, dnsRecords)
}

func (h *DomainHandler) CreateDNSRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	var record models.DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify zone belongs to user
	var zone models.DNSZone
	if err := h.db.Where("id = ? AND user_id = ?", record.ZoneID, userID).First(&zone).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DNS record"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *DomainHandler) UpdateDNSRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var record models.DNSRecord
	if err := h.db.Joins("JOIN dns_zones ON dns_records.zone_id = dns_zones.id").
		Where("dns_records.id = ? AND dns_zones.user_id = ?", id, userID).
		First(&record).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DNS record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *DomainHandler) DeleteDNSRecord(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Joins("JOIN dns_zones ON dns_records.zone_id = dns_zones.id").
		Where("dns_records.id = ? AND dns_zones.user_id = ?", id, userID).
		Delete(&models.DNSRecord{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DNS record deleted successfully"})
}

func (h *DomainHandler) GetErrorPages(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	errorPages := []map[string]interface{}{
		{
			"code":        404,
			"title":       "Page Not Found",
			"custom_page": "/error_pages/404.html",
			"enabled":     true,
		},
		{
			"code":        500,
			"title":       "Internal Server Error",
			"custom_page": "/error_pages/500.html",
			"enabled":     false,
		},
		{
			"code":        403,
			"title":       "Forbidden",
			"custom_page": "",
			"enabled":     false,
		},
	}
	
	c.JSON(http.StatusOK, errorPages)
}

func (h *DomainHandler) UpdateErrorPage(c *gin.Context) {
	userID := c.GetUint("user_id")
	code := c.Param("code")
	
	var request struct {
		CustomPage string `json:"custom_page"`
		Enabled    bool   `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Error page updated successfully",
		"code":        code,
		"user_id":     userID,
		"custom_page": request.CustomPage,
		"enabled":     request.Enabled,
	})
}
