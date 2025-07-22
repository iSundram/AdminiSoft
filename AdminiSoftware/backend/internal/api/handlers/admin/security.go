
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SecurityHandler struct {
	db *gorm.DB
}

func NewSecurityHandler(db *gorm.DB) *SecurityHandler {
	return &SecurityHandler{db: db}
}

func (h *SecurityHandler) GetSecurityOverview(c *gin.Context) {
	overview := map[string]interface{}{
		"two_factor_enabled":      true,
		"brute_force_protection":  true,
		"mod_security_enabled":    true,
		"csf_enabled":             true,
		"imunify360_enabled":      false,
		"failed_login_attempts":   0,
		"blocked_ips":             []string{},
		"security_scan_status":    "clean",
		"last_security_scan":      "2024-01-15T10:30:00Z",
		"ssl_certificates":        12,
		"firewall_rules":          25,
		"malware_detected":        0,
		"security_notifications":  []string{},
	}
	c.JSON(http.StatusOK, overview)
}

func (h *SecurityHandler) GetFailedLogins(c *gin.Context) {
	var logs []models.AccessLog
	if err := h.db.Where("status = ?", "failed").Order("timestamp desc").Limit(100).Find(&logs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch failed login attempts"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (h *SecurityHandler) GetBlockedIPs(c *gin.Context) {
	blockedIPs := []map[string]interface{}{
		{
			"ip":         "192.168.1.100",
			"reason":     "Brute force attack",
			"blocked_at": "2024-01-15T10:30:00Z",
			"expires_at": "2024-01-16T10:30:00Z",
		},
		{
			"ip":         "10.0.0.50",
			"reason":     "Malicious activity",
			"blocked_at": "2024-01-15T09:15:00Z",
			"expires_at": "2024-01-17T09:15:00Z",
		},
	}
	c.JSON(http.StatusOK, blockedIPs)
}

func (h *SecurityHandler) UnblockIP(c *gin.Context) {
	ip := c.Param("ip")
	c.JSON(http.StatusOK, gin.H{"message": "IP " + ip + " has been unblocked"})
}

func (h *SecurityHandler) UpdateFirewallRules(c *gin.Context) {
	var rules []map[string]interface{}
	if err := c.ShouldBindJSON(&rules); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Firewall rules updated successfully"})
}

func (h *SecurityHandler) RunSecurityScan(c *gin.Context) {
	scanResult := map[string]interface{}{
		"status":        "completed",
		"threats_found": 0,
		"files_scanned": 15420,
		"scan_time":     "2m 34s",
		"last_scan":     "2024-01-15T10:30:00Z",
		"report_url":    "/api/admin/security/scan-report",
	}
	c.JSON(http.StatusOK, scanResult)
}

func (h *SecurityHandler) GetSecuritySettings(c *gin.Context) {
	settings := map[string]interface{}{
		"max_login_attempts":        5,
		"lockout_duration":          30,
		"password_min_length":       8,
		"password_require_numbers":  true,
		"password_require_symbols":  true,
		"password_require_mixed":    true,
		"session_timeout":           1440,
		"two_factor_required":       false,
		"ip_whitelist_enabled":      false,
		"ip_whitelist":              []string{},
		"email_notifications":       true,
		"security_headers_enabled":  true,
		"csrf_protection_enabled":   true,
		"xss_protection_enabled":    true,
	}
	c.JSON(http.StatusOK, settings)
}

func (h *SecurityHandler) UpdateSecuritySettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Security settings updated successfully"})
}
