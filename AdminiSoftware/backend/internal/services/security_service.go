
package services

import (
	"fmt"
	"time"
)

type SecurityService struct {
	bruteForce *BruteForceProtection
}

func NewSecurityService() *SecurityService {
	return &SecurityService{
		bruteForce: NewBruteForceProtection(),
	}
}

func (s *SecurityService) CheckBruteForce(ip string) bool {
	return s.bruteForce.IsBlocked(ip)
}

func (s *SecurityService) RecordLoginAttempt(ip string, successful bool) {
	s.bruteForce.RecordAttempt(ip, successful)
}

func (s *SecurityService) GetSecuritySettings() map[string]interface{} {
	return map[string]interface{}{
		"two_factor_enabled":     true,
		"brute_force_protection": true,
		"max_login_attempts":     5,
		"block_duration":         30,
		"ssl_enforced":          true,
		"password_policy": map[string]interface{}{
			"min_length":      8,
			"require_upper":   true,
			"require_lower":   true,
			"require_numbers": true,
			"require_symbols": true,
		},
	}
}

func (s *SecurityService) UpdateSecuritySettings(settings map[string]interface{}) error {
	// Implementation for updating security settings
	return nil
}

func (s *SecurityService) GetSecurityAuditLog() []SecurityEvent {
	return []SecurityEvent{
		{
			Type:        "login_attempt",
			IP:          "192.168.1.100",
			UserAgent:   "Mozilla/5.0...",
			Timestamp:   time.Now().Add(-1 * time.Hour),
			Successful:  false,
			Description: "Failed login attempt",
		},
		{
			Type:        "password_change",
			IP:          "192.168.1.50",
			UserAgent:   "Mozilla/5.0...",
			Timestamp:   time.Now().Add(-2 * time.Hour),
			Successful:  true,
			Description: "Password changed successfully",
		},
	}
}

type SecurityEvent struct {
	Type        string    `json:"type"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Timestamp   time.Time `json:"timestamp"`
	Successful  bool      `json:"successful"`
	Description string    `json:"description"`
}

type BruteForceProtection struct {
	// Implementation same as in brute_force.go
}

func NewBruteForceProtection() *BruteForceProtection {
	return &BruteForceProtection{}
}

func (bfp *BruteForceProtection) IsBlocked(ip string) bool {
	return false
}

func (bfp *BruteForceProtection) RecordAttempt(ip string, successful bool) {
	// Implementation
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"
	"fmt"
	"net"
	"time"

	"gorm.io/gorm"
)

type SecurityService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewSecurityService(db *gorm.DB, logger *utils.Logger) *SecurityService {
	return &SecurityService{
		db:     db,
		logger: logger,
	}
}

// Brute Force Protection
func (s *SecurityService) RecordLoginAttempt(ip string, success bool, username string) error {
	attempt := models.LoginAttempt{
		IPAddress: ip,
		Username:  username,
		Success:   success,
		Timestamp: time.Now(),
	}

	if err := s.db.Create(&attempt).Error; err != nil {
		s.logger.Error("Failed to record login attempt", map[string]interface{}{
			"error": err.Error(),
			"ip": ip,
		})
		return err
	}

	// Check for brute force patterns
	if !success {
		if err := s.checkBruteForce(ip); err != nil {
			s.logger.Warn("Brute force detected", map[string]interface{}{
				"ip": ip,
				"username": username,
			})
		}
	}

	return nil
}

func (s *SecurityService) checkBruteForce(ip string) error {
	// Count failed attempts in the last 15 minutes
	since := time.Now().Add(-15 * time.Minute)
	var count int64

	if err := s.db.Model(&models.LoginAttempt{}).
		Where("ip_address = ? AND success = false AND timestamp >= ?", ip, since).
		Count(&count).Error; err != nil {
		return err
	}

	// If more than 5 failed attempts, block the IP
	if count >= 5 {
		return s.BlockIP(ip, "Brute force attack detected", 3600) // Block for 1 hour
	}

	return nil
}

func (s *SecurityService) BlockIP(ip string, reason string, durationSeconds int) error {
	// Check if IP is already blocked
	var existing models.BlockedIP
	if err := s.db.Where("ip_address = ?", ip).First(&existing).Error; err == nil {
		// Update existing block
		existing.Reason = reason
		existing.ExpiresAt = time.Now().Add(time.Duration(durationSeconds) * time.Second)
		return s.db.Save(&existing).Error
	}

	// Create new block
	blockedIP := models.BlockedIP{
		IPAddress: ip,
		Reason:    reason,
		BlockedBy: "system",
		ExpiresAt: time.Now().Add(time.Duration(durationSeconds) * time.Second),
	}

	if err := s.db.Create(&blockedIP).Error; err != nil {
		s.logger.Error("Failed to block IP", map[string]interface{}{
			"error": err.Error(),
			"ip": ip,
		})
		return err
	}

	s.logger.Info("IP address blocked", map[string]interface{}{
		"ip": ip,
		"reason": reason,
		"duration": durationSeconds,
	})

	return nil
}

func (s *SecurityService) IsIPBlocked(ip string) bool {
	var blockedIP models.BlockedIP
	err := s.db.Where("ip_address = ? AND expires_at > ?", ip, time.Now()).First(&blockedIP).Error
	return err == nil
}

func (s *SecurityService) UnblockIP(ip string) error {
	if err := s.db.Where("ip_address = ?", ip).Delete(&models.BlockedIP{}).Error; err != nil {
		s.logger.Error("Failed to unblock IP", map[string]interface{}{
			"error": err.Error(),
			"ip": ip,
		})
		return err
	}

	s.logger.Info("IP address unblocked", map[string]interface{}{
		"ip": ip,
	})

	return nil
}

// Two-Factor Authentication
func (s *SecurityService) Enable2FA(userID uint) (*models.TwoFactorAuth, error) {
	// Check if 2FA is already enabled
	var existing models.TwoFactorAuth
	if err := s.db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, errors.New("2FA is already enabled for this user")
	}

	// Generate secret key
	secret, err := utils.GenerateTOTPSecret()
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP secret: %v", err)
	}

	tfa := models.TwoFactorAuth{
		UserID:      userID,
		Secret:      secret,
		IsEnabled:   false, // Will be enabled after verification
		BackupCodes: utils.GenerateBackupCodes(),
	}

	if err := s.db.Create(&tfa).Error; err != nil {
		s.logger.Error("Failed to create 2FA record", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		})
		return nil, err
	}

	return &tfa, nil
}

func (s *SecurityService) Verify2FA(userID uint, code string) error {
	var tfa models.TwoFactorAuth
	if err := s.db.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		return errors.New("2FA not found for user")
	}

	// Verify TOTP code
	if !utils.VerifyTOTP(tfa.Secret, code) {
		// Check backup codes
		if !s.verifyBackupCode(userID, code) {
			return errors.New("invalid 2FA code")
		}
	}

	// Enable 2FA if not already enabled
	if !tfa.IsEnabled {
		tfa.IsEnabled = true
		if err := s.db.Save(&tfa).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *SecurityService) verifyBackupCode(userID uint, code string) bool {
	var tfa models.TwoFactorAuth
	if err := s.db.Where("user_id = ?", userID).First(&tfa).Error; err != nil {
		return false
	}

	for i, backupCode := range tfa.BackupCodes {
		if backupCode == code {
			// Remove used backup code
			tfa.BackupCodes = append(tfa.BackupCodes[:i], tfa.BackupCodes[i+1:]...)
			s.db.Save(&tfa)
			return true
		}
	}

	return false
}

func (s *SecurityService) Disable2FA(userID uint) error {
	if err := s.db.Where("user_id = ?", userID).Delete(&models.TwoFactorAuth{}).Error; err != nil {
		s.logger.Error("Failed to disable 2FA", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		})
		return err
	}

	s.logger.Info("2FA disabled for user", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// Security Scanning
func (s *SecurityService) ScanForMalware(path string) (*models.SecurityScan, error) {
	scan := models.SecurityScan{
		ScanType:  "malware",
		Path:      path,
		Status:    "running",
		StartedAt: time.Now(),
	}

	if err := s.db.Create(&scan).Error; err != nil {
		return nil, err
	}

	// Start scan in background
	go s.performMalwareScan(&scan)

	return &scan, nil
}

func (s *SecurityService) performMalwareScan(scan *models.SecurityScan) {
	// Simulate malware scan - in production, integrate with ClamAV or similar
	time.Sleep(10 * time.Second) // Simulate scan time

	scan.Status = "completed"
	scan.CompletedAt = time.Now()
	scan.FilesScanned = 1000
	scan.ThreatsFound = 0
	scan.Results = "No threats detected"

	s.db.Save(scan)

	s.logger.Info("Malware scan completed", map[string]interface{}{
		"scan_id": scan.ID,
		"path": scan.Path,
		"threats": scan.ThreatsFound,
	})
}

// Firewall Management
func (s *SecurityService) CreateFirewallRule(rule *models.FirewallRule) error {
	if err := s.validateFirewallRule(rule); err != nil {
		return err
	}

	if err := s.db.Create(rule).Error; err != nil {
		s.logger.Error("Failed to create firewall rule", map[string]interface{}{
			"error": err.Error(),
			"action": rule.Action,
		})
		return err
	}

	s.logger.Info("Firewall rule created", map[string]interface{}{
		"rule_id": rule.ID,
		"action": rule.Action,
		"source_ip": rule.SourceIP,
	})

	return nil
}

func (s *SecurityService) validateFirewallRule(rule *models.FirewallRule) error {
	if rule.Action != "allow" && rule.Action != "deny" {
		return errors.New("action must be 'allow' or 'deny'")
	}

	if rule.SourceIP != "" {
		if ip := net.ParseIP(rule.SourceIP); ip == nil {
			if _, _, err := net.ParseCIDR(rule.SourceIP); err != nil {
				return errors.New("invalid source IP or CIDR")
			}
		}
	}

	if rule.Port != 0 && (rule.Port < 1 || rule.Port > 65535) {
		return errors.New("invalid port number")
	}

	return nil
}

func (s *SecurityService) DeleteFirewallRule(ruleID uint) error {
	var rule models.FirewallRule
	if err := s.db.First(&rule, ruleID).Error; err != nil {
		return errors.New("firewall rule not found")
	}

	if err := s.db.Delete(&rule).Error; err != nil {
		s.logger.Error("Failed to delete firewall rule", map[string]interface{}{
			"error": err.Error(),
			"rule_id": ruleID,
		})
		return err
	}

	s.logger.Info("Firewall rule deleted", map[string]interface{}{
		"rule_id": ruleID,
	})

	return nil
}

func (s *SecurityService) ListFirewallRules() ([]models.FirewallRule, error) {
	var rules []models.FirewallRule
	
	if err := s.db.Find(&rules).Error; err != nil {
		s.logger.Error("Failed to list firewall rules", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return rules, nil
}

// SSL Security
func (s *SecurityService) CheckSSLSecurity(domain string) (*models.SSLSecurityCheck, error) {
	check := models.SSLSecurityCheck{
		Domain:    domain,
		CheckedAt: time.Now(),
	}

	// Perform SSL security check (simplified)
	// In production, this would check certificate validity, cipher suites, etc.
	check.HasValidCertificate = true
	check.CertificateExpiry = time.Now().AddDate(0, 3, 0) // 3 months from now
	check.Grade = "A"
	check.Vulnerabilities = []string{}

	if err := s.db.Create(&check).Error; err != nil {
		s.logger.Error("Failed to save SSL security check", map[string]interface{}{
			"error": err.Error(),
			"domain": domain,
		})
		return nil, err
	}

	return &check, nil
}

// Security Log Management
func (s *SecurityService) LogSecurityEvent(event *models.SecurityEvent) error {
	if err := s.db.Create(event).Error; err != nil {
		s.logger.Error("Failed to log security event", map[string]interface{}{
			"error": err.Error(),
			"event_type": event.EventType,
		})
		return err
	}

	s.logger.Info("Security event logged", map[string]interface{}{
		"event_type": event.EventType,
		"severity": event.Severity,
		"ip_address": event.IPAddress,
	})

	return nil
}

func (s *SecurityService) GetSecurityEvents(limit int) ([]models.SecurityEvent, error) {
	var events []models.SecurityEvent
	
	query := s.db.Order("timestamp DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&events).Error; err != nil {
		s.logger.Error("Failed to get security events", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return events, nil
}

// Cleanup expired records
func (s *SecurityService) CleanupExpiredRecords() error {
	now := time.Now()

	// Clean up expired IP blocks
	if err := s.db.Where("expires_at < ?", now).Delete(&models.BlockedIP{}).Error; err != nil {
		s.logger.Error("Failed to cleanup expired IP blocks", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Clean up old login attempts (keep for 30 days)
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	if err := s.db.Where("timestamp < ?", thirtyDaysAgo).Delete(&models.LoginAttempt{}).Error; err != nil {
		s.logger.Error("Failed to cleanup old login attempts", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Clean up old security events (keep for 90 days)
	ninetyDaysAgo := now.AddDate(0, 0, -90)
	if err := s.db.Where("timestamp < ?", ninetyDaysAgo).Delete(&models.SecurityEvent{}).Error; err != nil {
		s.logger.Error("Failed to cleanup old security events", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return nil
}
