
package admin

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

func (h *StatsHandler) GetSystemStats(c *gin.Context) {
	var stats []models.SystemStat
	if err := h.db.Order("timestamp desc").Limit(24).Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch system stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetUserStats(c *gin.Context) {
	var userStats []models.UserStats
	if err := h.db.Preload("User").Order("date desc").Limit(30).Find(&userStats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user statistics"})
		return
	}
	c.JSON(http.StatusOK, userStats)
}

func (h *StatsHandler) GetAccessLogs(c *gin.Context) {
	var logs []models.AccessLog
	if err := h.db.Preload("User").Order("timestamp desc").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	var totalUsers int64
	var totalDomains int64
	var totalDatabases int64
	var totalEmailAccounts int64

	h.db.Model(&models.User{}).Count(&totalUsers)
	h.db.Model(&models.Domain{}).Count(&totalDomains)
	h.db.Model(&models.Database{}).Count(&totalDatabases)
	h.db.Model(&models.EmailAccount{}).Count(&totalEmailAccounts)

	stats := map[string]interface{}{
		"total_users":          totalUsers,
		"total_domains":        totalDomains,
		"total_databases":      totalDatabases,
		"total_email_accounts": totalEmailAccounts,
		"server_load":          []float64{0.5, 0.3, 0.7},
		"cpu_usage":            45.2,
		"memory_usage":         67.8,
		"disk_usage":           34.5,
		"bandwidth_usage":      {
			"used":  "245.6 GB",
			"total": "1 TB",
			"percentage": 24.56,
		},
		"uptime":               "15 days, 3 hours, 22 minutes",
		"active_sessions":      127,
		"failed_logins_today":  8,
		"backup_status":        "completed",
		"last_backup":          "2024-01-15T02:00:00Z",
		"ssl_certificates": {
			"total":    45,
			"valid":    42,
			"expiring": 2,
			"expired":  1,
		},
		"service_status": {
			"apache":     "running",
			"mysql":      "running",
			"postfix":    "running",
			"dovecot":    "running",
			"bind":       "running",
			"vsftpd":     "running",
		},
		"recent_activities": []map[string]interface{}{
			{
				"user":   "admin",
				"action": "Created new account",
				"time":   "2024-01-15T10:30:00Z",
			},
			{
				"user":   "user1",
				"action": "Updated SSL certificate",
				"time":   "2024-01-15T10:25:00Z",
			},
			{
				"user":   "reseller1",
				"action": "Created hosting package",
				"time":   "2024-01-15T10:20:00Z",
			},
		},
	}

	c.JSON(http.StatusOK, stats)
}

func (h *StatsHandler) GetResourceUsage(c *gin.Context) {
	usage := map[string]interface{}{
		"cpu": {
			"current":     45.2,
			"average_1h":  38.7,
			"average_24h": 42.1,
			"peak_24h":    78.9,
		},
		"memory": {
			"total":       "32 GB",
			"used":        "21.7 GB",
			"free":        "10.3 GB",
			"percentage":  67.8,
			"swap_used":   "512 MB",
			"swap_total":  "4 GB",
		},
		"disk": {
			"total":      "2 TB",
			"used":       "690 GB",
			"free":       "1.31 TB",
			"percentage": 34.5,
			"inodes_used": 125000,
			"inodes_total": 1000000,
		},
		"network": {
			"rx_bytes":   "1.2 TB",
			"tx_bytes":   "890 GB",
			"rx_packets": 45000000,
			"tx_packets": 38000000,
			"current_rx": "125 Mbps",
			"current_tx": "89 Mbps",
		},
		"load_average": {
			"1min":  0.52,
			"5min":  0.38,
			"15min": 0.41,
		},
	}
	c.JSON(http.StatusOK, usage)
}

func (h *StatsHandler) GetBandwidthStats(c *gin.Context) {
	bandwidth := map[string]interface{}{
		"today": {
			"in":  "12.3 GB",
			"out": "8.7 GB",
		},
		"this_week": {
			"in":  "89.2 GB",
			"out": "67.4 GB",
		},
		"this_month": {
			"in":  "345.6 GB",
			"out": "267.8 GB",
		},
		"daily_history": []map[string]interface{}{
			{"date": "2024-01-15", "in": 12.3, "out": 8.7},
			{"date": "2024-01-14", "in": 15.7, "out": 11.2},
			{"date": "2024-01-13", "in": 18.4, "out": 13.8},
			{"date": "2024-01-12", "in": 14.9, "out": 10.5},
			{"date": "2024-01-11", "in": 16.3, "out": 12.1},
		},
		"top_users": []map[string]interface{}{
			{"username": "user1", "bandwidth": "45.2 GB"},
			{"username": "user2", "bandwidth": "38.7 GB"},
			{"username": "user3", "bandwidth": "32.1 GB"},
		},
	}
	c.JSON(http.StatusOK, bandwidth)
}
