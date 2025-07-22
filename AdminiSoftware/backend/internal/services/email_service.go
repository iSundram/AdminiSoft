
package services

import (
	"AdminiSoftware/internal/models"
	"gorm.io/gorm"
)

type EmailService struct {
	db *gorm.DB
}

func NewEmailService(db *gorm.DB) *EmailService {
	return &EmailService{db: db}
}

func (s *EmailService) CreateEmailAccount(account *models.EmailAccount) error {
	return s.db.Create(account).Error
}

func (s *EmailService) GetEmailAccounts(userID uint) ([]models.EmailAccount, error) {
	var accounts []models.EmailAccount
	err := s.db.Where("user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (s *EmailService) UpdateEmailAccount(account *models.EmailAccount) error {
	return s.db.Save(account).Error
}

func (s *EmailService) DeleteEmailAccount(id uint) error {
	return s.db.Delete(&models.EmailAccount{}, id).Error
}

func (s *EmailService) CreateForwarder(forwarder *models.EmailForwarder) error {
	return s.db.Create(forwarder).Error
}

func (s *EmailService) GetForwarders(userID uint) ([]models.EmailForwarder, error) {
	var forwarders []models.EmailForwarder
	err := s.db.Where("user_id = ?", userID).Find(&forwarders).Error
	return forwarders, err
}

func (s *EmailService) DeleteForwarder(id uint) error {
	return s.db.Delete(&models.EmailForwarder{}, id).Error
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type EmailService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewEmailService(db *gorm.DB, logger *utils.Logger) *EmailService {
	return &EmailService{
		db:     db,
		logger: logger,
	}
}

func (s *EmailService) CreateEmailAccount(email *models.EmailAccount) error {
	// Validate email
	if email.Email == "" || email.Domain == "" {
		return errors.New("email and domain are required")
	}

	// Check if email already exists
	var existing models.EmailAccount
	if err := s.db.Where("email = ?", email.Email).First(&existing).Error; err == nil {
		return errors.New("email account already exists")
	}

	// Hash password
	if email.Password != "" {
		hashedPassword, err := utils.HashPassword(email.Password)
		if err != nil {
			return err
		}
		email.PasswordHash = hashedPassword
		email.Password = "" // Clear plaintext password
	}

	// Set default quotas if not specified
	if email.QuotaMB == 0 {
		email.QuotaMB = 1000 // 1GB default
	}

	if err := s.db.Create(email).Error; err != nil {
		s.logger.Error("Failed to create email account", map[string]interface{}{
			"error": err.Error(),
			"email": email.Email,
		})
		return err
	}

	return nil
}

func (s *EmailService) UpdateEmailAccount(emailID uint, updates map[string]interface{}) error {
	var email models.EmailAccount
	if err := s.db.First(&email, emailID).Error; err != nil {
		return errors.New("email account not found")
	}

	// Hash password if provided
	if password, ok := updates["password"]; ok && password != "" {
		hashedPassword, err := utils.HashPassword(password.(string))
		if err != nil {
			return err
		}
		updates["password_hash"] = hashedPassword
		delete(updates, "password")
	}

	if err := s.db.Model(&email).Updates(updates).Error; err != nil {
		s.logger.Error("Failed to update email account", map[string]interface{}{
			"error": err.Error(),
			"email_id": emailID,
		})
		return err
	}

	return nil
}

func (s *EmailService) DeleteEmailAccount(emailID uint) error {
	var email models.EmailAccount
	if err := s.db.First(&email, emailID).Error; err != nil {
		return errors.New("email account not found")
	}

	if err := s.db.Delete(&email).Error; err != nil {
		s.logger.Error("Failed to delete email account", map[string]interface{}{
			"error": err.Error(),
			"email_id": emailID,
		})
		return err
	}

	return nil
}

func (s *EmailService) ListEmailAccounts(userID uint, role string) ([]models.EmailAccount, error) {
	var emails []models.EmailAccount
	query := s.db

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&emails).Error; err != nil {
		s.logger.Error("Failed to list email accounts", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		})
		return nil, err
	}

	return emails, nil
}

func (s *EmailService) CreateForwarder(forwarder *models.EmailForwarder) error {
	if forwarder.SourceEmail == "" || forwarder.DestinationEmail == "" {
		return errors.New("source and destination emails are required")
	}

	if err := s.db.Create(forwarder).Error; err != nil {
		s.logger.Error("Failed to create email forwarder", map[string]interface{}{
			"error": err.Error(),
			"source": forwarder.SourceEmail,
		})
		return err
	}

	return nil
}

func (s *EmailService) DeleteForwarder(forwarderID uint) error {
	var forwarder models.EmailForwarder
	if err := s.db.First(&forwarder, forwarderID).Error; err != nil {
		return errors.New("email forwarder not found")
	}

	if err := s.db.Delete(&forwarder).Error; err != nil {
		s.logger.Error("Failed to delete email forwarder", map[string]interface{}{
			"error": err.Error(),
			"forwarder_id": forwarderID,
		})
		return err
	}

	return nil
}

func (s *EmailService) GetEmailStats(domain string) (*models.EmailStats, error) {
	var stats models.EmailStats
	
	// Count email accounts
	if err := s.db.Model(&models.EmailAccount{}).Where("domain = ?", domain).Count(&stats.TotalAccounts).Error; err != nil {
		return nil, err
	}

	// Count forwarders
	if err := s.db.Model(&models.EmailForwarder{}).Where("domain = ?", domain).Count(&stats.TotalForwarders).Error; err != nil {
		return nil, err
	}

	// Calculate total quota used
	var totalUsed int64
	if err := s.db.Model(&models.EmailAccount{}).Where("domain = ?", domain).Select("COALESCE(SUM(used_mb), 0)").Scan(&totalUsed).Error; err != nil {
		return nil, err
	}
	stats.QuotaUsedMB = totalUsed

	return &stats, nil
}

func (s *EmailService) ChangeEmailPassword(emailID uint, newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	var email models.EmailAccount
	if err := s.db.First(&email, emailID).Error; err != nil {
		return errors.New("email account not found")
	}

	if err := s.db.Model(&email).Update("password_hash", hashedPassword).Error; err != nil {
		s.logger.Error("Failed to change email password", map[string]interface{}{
			"error": err.Error(),
			"email_id": emailID,
		})
		return err
	}

	return nil
}
