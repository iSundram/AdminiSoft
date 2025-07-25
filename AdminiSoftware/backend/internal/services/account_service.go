
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AccountService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewAccountService(db *gorm.DB, logger *utils.Logger) *AccountService {
	return &AccountService{
		db:     db,
		logger: logger,
	}
}

type AccountStats struct {
	TotalAccounts    int64 `json:"total_accounts"`
	ActiveAccounts   int64 `json:"active_accounts"`
	SuspendedAccounts int64 `json:"suspended_accounts"`
	OverQuotaAccounts int64 `json:"over_quota_accounts"`
}

func (s *AccountService) ListAccounts(page, limit int, search, status string) ([]models.User, int64, error) {
	var accounts []models.User
	var total int64

	query := s.db.Model(&models.User{}).Where("role != ?", "admin")

	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("Package").Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (s *AccountService) CreateAccount(username, email, password, domain string, packageID, resellerID uint) (*models.User, error) {
	// Check if username already exists
	var existingUser models.User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	}

	// Validate package exists
	var pkg models.Package
	if err := s.db.First(&pkg, packageID).Error; err != nil {
		return nil, errors.New("package not found")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user account
	user := models.User{
		Username:    username,
		Email:       email,
		Password:    hashedPassword,
		Role:        "user",
		Status:      "active",
		PackageID:   packageID,
		ResellerID:  resellerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Create domain entry
	domainEntry := models.Domain{
		Name:      domain,
		UserID:    user.ID,
		Type:      "main",
		Status:    "active",
		CreatedAt: time.Now(),
	}

	if err := s.db.Create(&domainEntry).Error; err != nil {
		// Rollback user creation if domain creation fails
		s.db.Delete(&user)
		return nil, err
	}

	s.logger.Info(fmt.Sprintf("Account created: %s (%s)", username, email))
	return &user, nil
}

func (s *AccountService) GetAccount(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Package").Preload("Domains").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AccountService) UpdateAccount(id uint, email string, packageID uint, status string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if email != "" {
		user.Email = email
	}
	if packageID != 0 {
		user.PackageID = packageID
	}
	if status != "" {
		user.Status = status
	}

	user.UpdatedAt = time.Now()

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AccountService) SuspendAccount(id uint, reason string) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	user.Status = "suspended"
	user.SuspensionReason = reason
	user.UpdatedAt = time.Now()

	return s.db.Save(&user).Error
}

func (s *AccountService) UnsuspendAccount(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	user.Status = "active"
	user.SuspensionReason = ""
	user.UpdatedAt = time.Now()

	return s.db.Save(&user).Error
}

func (s *AccountService) DeleteAccount(id uint) error {
	// Start transaction
	tx := s.db.Begin()

	// Delete related domains
	if err := tx.Where("user_id = ?", id).Delete(&models.Domain{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete user
	if err := tx.Delete(&models.User{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *AccountService) GetAccountStats() (*AccountStats, error) {
	var stats AccountStats

	// Total accounts (excluding admins)
	s.db.Model(&models.User{}).Where("role != ?", "admin").Count(&stats.TotalAccounts)
	
	// Active accounts
	s.db.Model(&models.User{}).Where("role != ? AND status = ?", "admin", "active").Count(&stats.ActiveAccounts)
	
	// Suspended accounts
	s.db.Model(&models.User{}).Where("role != ? AND status = ?", "admin", "suspended").Count(&stats.SuspendedAccounts)
	
	// Over quota accounts (this would need resource tracking)
	stats.OverQuotaAccounts = 0

	return &stats, nil
}
