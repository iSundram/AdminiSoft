
package reseller

import (
	"AdminiSoftware/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	
	var totalAccounts int64
	var activeAccounts int64
	var totalPackages int64
	var totalDomains int64

	h.db.Model(&models.User{}).Where("reseller_id = ?", resellerID).Count(&totalAccounts)
	h.db.Model(&models.User{}).Where("reseller_id = ? AND status = ?", resellerID, "active").Count(&activeAccounts)
	h.db.Model(&models.Package{}).Where("reseller_id = ?", resellerID).Count(&totalPackages)
	h.db.Model(&models.Domain{}).Joins("JOIN users ON domains.user_id = users.id").Where("users.reseller_id = ?", resellerID).Count(&totalDomains)

	stats := map[string]interface{}{
		"total_accounts":     totalAccounts,
		"active_accounts":    activeAccounts,
		"suspended_accounts": totalAccounts - activeAccounts,
		"total_packages":     totalPackages,
		"total_domains":      totalDomains,
		"total_bandwidth":    "2.5 TB",
		"used_bandwidth":     "1.8 TB",
		"total_disk_space":   "500 GB",
		"used_disk_space":    "342 GB",
		"monthly_revenue":    "$2,340",
		"new_signups_today":  3,
		"new_signups_week":   18,
		"new_signups_month":  76,
		"resource_usage": map[string]interface{}{
			"disk_percentage":      68.4,
			"bandwidth_percentage": 72.0,
			"accounts_percentage":  45.5,
		},
		"recent_activities": []map[string]interface{}{
			{
				"user":   "newuser1",
				"action": "Account created",
				"time":   "2024-01-15T10:30:00Z",
			},
			{
				"user":   "customer1",
				"action": "Domain added",
				"time":   "2024-01-15T10:25:00Z",
			},
			{
				"user":   "customer2",
				"action": "Package upgraded",
				"time":   "2024-01-15T10:20:00Z",
			},
		},
		"top_customers": []map[string]interface{}{
			{
				"username":   "customer1",
				"disk_used":  "15.6 GB",
				"bandwidth":  "145.2 GB",
			},
			{
				"username":   "customer2",
				"disk_used":  "12.3 GB",
				"bandwidth":  "98.7 GB",
			},
			{
				"username":   "customer3",
				"disk_used":  "10.8 GB",
				"bandwidth":  "76.4 GB",
			},
		},
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetResourceUsage(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	
	usage := map[string]interface{}{
		"reseller_id": resellerID,
		"limits": map[string]interface{}{
			"accounts":   100,
			"disk_space": "500 GB",
			"bandwidth":  "2.5 TB",
			"domains":    500,
		},
		"used": map[string]interface{}{
			"accounts":   45,
			"disk_space": "342 GB",
			"bandwidth":  "1.8 TB",
			"domains":    234,
		},
		"percentages": map[string]interface{}{
			"accounts":   45.0,
			"disk_space": 68.4,
			"bandwidth":  72.0,
			"domains":    46.8,
		},
		"monthly_usage": []map[string]interface{}{
			{"month": "Jan", "disk": 342, "bandwidth": 1800, "accounts": 45},
			{"month": "Dec", "disk": 298, "bandwidth": 1650, "accounts": 42},
			{"month": "Nov", "disk": 267, "bandwidth": 1420, "accounts": 38},
		},
	}

	c.JSON(http.StatusOK, usage)
}

func (h *StatsHandler) GetBandwidthStats(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	
	bandwidth := map[string]interface{}{
		"reseller_id": resellerID,
		"current_month": map[string]interface{}{
			"used":  "1.8 TB",
			"limit": "2.5 TB",
			"percentage": 72.0,
		},
		"daily_usage": []map[string]interface{}{
			{"date": "2024-01-15", "usage": 65.2},
			{"date": "2024-01-14", "usage": 58.7},
			{"date": "2024-01-13", "usage": 72.4},
			{"date": "2024-01-12", "usage": 69.1},
			{"date": "2024-01-11", "usage": 54.8},
		},
		"top_accounts": []map[string]interface{}{
			{"username": "customer1", "usage": "145.2 GB"},
			{"username": "customer2", "usage": "98.7 GB"},
			{"username": "customer3", "usage": "76.4 GB"},
		},
		"total_this_year": "18.7 TB",
		"average_daily":   "62.3 GB",
		"peak_usage":      "89.4 GB",
	}

	c.JSON(http.StatusOK, bandwidth)
}
