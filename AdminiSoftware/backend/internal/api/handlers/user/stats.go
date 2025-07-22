
package user

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserStatsHandler struct {
	db *gorm.DB
}

func NewUserStatsHandler(db *gorm.DB) *UserStatsHandler {
	return &UserStatsHandler{db: db}
}

func (h *UserStatsHandler) GetDashboardStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get domain count
	var domainCount int64
	h.db.Model(&models.Domain{}).Where("user_id = ?", userID).Count(&domainCount)

	// Get email account count
	var emailCount int64
	h.db.Model(&models.Email{}).Where("user_id = ?", userID).Count(&emailCount)

	// Get database count
	var dbCount int64
	h.db.Model(&models.Database{}).Where("user_id = ?", userID).Count(&dbCount)

	// Get SSL certificate count
	var sslCount int64
	h.db.Model(&models.SSL{}).Where("user_id = ? AND status = 'active'", userID).Count(&sslCount)

	// Get disk usage
	var user models.User
	h.db.First(&user, userID)

	stats := map[string]interface{}{
		"domains":      domainCount,
		"emails":       emailCount,
		"databases":    dbCount,
		"ssl_certs":    sslCount,
		"disk_used":    user.DiskUsed,
		"disk_limit":   user.DiskLimit,
		"bandwidth_used": user.BandwidthUsed,
		"bandwidth_limit": user.BandwidthLimit,
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

func (h *UserStatsHandler) GetBandwidthUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var stats []models.BandwidthStat
	if err := h.db.Where("user_id = ? AND date >= ?", userID, time.Now().AddDate(0, -1, 0)).Order("date ASC").Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bandwidth usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"bandwidth_stats": stats})
}

func (h *UserStatsHandler) GetDiskUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	type DiskUsageBreakdown struct {
		Category string  `json:"category"`
		Size     int64   `json:"size"`
		Percentage float64 `json:"percentage"`
	}

	var user models.User
	h.db.First(&user, userID)

	// Simulate disk usage breakdown
	breakdown := []DiskUsageBreakdown{
		{Category: "Web Files", Size: user.DiskUsed * 40 / 100, Percentage: 40},
		{Category: "Email", Size: user.DiskUsed * 25 / 100, Percentage: 25},
		{Category: "Databases", Size: user.DiskUsed * 20 / 100, Percentage: 20},
		{Category: "Backups", Size: user.DiskUsed * 10 / 100, Percentage: 10},
		{Category: "Other", Size: user.DiskUsed * 5 / 100, Percentage: 5},
	}

	c.JSON(http.StatusOK, gin.H{
		"total_used": user.DiskUsed,
		"total_limit": user.DiskLimit,
		"breakdown": breakdown,
	})
}

func (h *UserStatsHandler) GetVisitorStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	domainID := c.Query("domain_id")

	var stats []models.VisitorStat
	query := h.db.Where("user_id = ?", userID)
	
	if domainID != "" {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Where("date >= ?", time.Now().AddDate(0, 0, -30)).Order("date ASC").Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visitor statistics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"visitor_stats": stats})
}

func (h *UserStatsHandler) GetErrorLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	domainID := c.Query("domain_id")
	limit := c.DefaultQuery("limit", "100")

	var logs []models.ErrorLog
	query := h.db.Where("user_id = ?", userID)
	
	if domainID != "" {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Order("created_at DESC").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch error logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error_logs": logs, "limit": limit})
}

func (h *UserStatsHandler) GetAccessLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	domainID := c.Query("domain_id")
	limit := c.DefaultQuery("limit", "1000")

	var logs []models.AccessLog
	query := h.db.Where("user_id = ?", userID)
	
	if domainID != "" {
		query = query.Where("domain_id = ?", domainID)
	}

	if err := query.Order("timestamp DESC").Limit(1000).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_logs": logs, "limit": limit})
}

func (h *UserStatsHandler) GetResourceUsage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var usage []models.ResourceUsage
	if err := h.db.Where("user_id = ? AND timestamp >= ?", userID, time.Now().Add(-24*time.Hour)).Order("timestamp ASC").Find(&usage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resource usage"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"resource_usage": usage})
}
