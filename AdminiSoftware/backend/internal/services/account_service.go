
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
package services

import (
	"AdminiSoftware/internal/models"
	"fmt"
	"gorm.io/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) CreateAccount(account *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	account.Password = string(hashedPassword)

	return s.db.Create(account).Error
}

func (s *AccountService) GetAccountByID(id uint) (*models.User, error) {
	var account models.User
	err := s.db.Preload("Package").First(&account, id).Error
	return &account, err
}

func (s *AccountService) GetAccountsByReseller(resellerID uint) ([]models.User, error) {
	var accounts []models.User
	err := s.db.Where("reseller_id = ?", resellerID).Preload("Package").Find(&accounts).Error
	return accounts, err
}

func (s *AccountService) UpdateAccount(account *models.User) error {
	return s.db.Save(account).Error
}

func (s *AccountService) SuspendAccount(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", "suspended").Error
}

func (s *AccountService) UnsuspendAccount(id uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", id).Update("status", "active").Error
}

func (s *AccountService) DeleteAccount(id uint) error {
	return s.db.Delete(&models.User{}, id).Error
}

func (s *AccountService) GetAccountStats(id uint) (map[string]interface{}, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"disk_used":      user.DiskUsed,
		"bandwidth_used": user.BandwidthUsed,
		"domains_count":  s.countUserDomains(id),
		"email_accounts": s.countUserEmails(id),
		"databases":      s.countUserDatabases(id),
	}

	return stats, nil
}

func (s *AccountService) countUserDomains(userID uint) int64 {
	var count int64
	s.db.Model(&models.Domain{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (s *AccountService) countUserEmails(userID uint) int64 {
	var count int64
	s.db.Model(&models.Email{}).Where("user_id = ?", userID).Count(&count)
	return count
}

func (s *AccountService) countUserDatabases(userID uint) int64 {
	var count int64
	s.db.Model(&models.Database{}).Where("user_id = ?", userID).Count(&count)
	return count
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"
	"fmt"

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

func (s *AccountService) GetAccountsPaginated(page, limit int) ([]models.User, int64, error) {
	var accounts []models.User
	var total int64

	offset := (page - 1) * limit

	if err := s.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Preload("Package").Offset(offset).Limit(limit).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

func (s *AccountService) CreateAccount(req *models.CreateAccountRequest) (*models.User, error) {
	// Check if username already exists
	var existing models.User
	if err := s.db.Where("username = ?", req.Username).First(&existing).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if err := s.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		PackageID:    req.PackageID,
		Role:         "user",
		Status:       "active",
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Create domain
	domain := &models.Domain{
		UserID:       user.ID,
		Domain:       req.Domain,
		DocumentRoot: "/public_html",
		Status:       "active",
	}

	if err := s.db.Create(domain).Error; err != nil {
		s.logger.Error("Failed to create domain for user: " + err.Error())
		// Don't fail the account creation for domain creation failure
	}

	// Load package information
	if err := s.db.Preload("Package").First(user, user.ID).Error; err != nil {
		s.logger.Error("Failed to load package info: " + err.Error())
	}

	s.logger.Info(fmt.Sprintf("Account created successfully for user: %s", user.Username))
	return user, nil
}

func (s *AccountService) UpdateAccount(id uint, req *models.UpdateAccountRequest) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, errors.New("account not found")
	}

	// Update fields if provided
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.PackageID != 0 {
		user.PackageID = req.PackageID
	}
	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update account: %v", err)
	}

	// Load package information
	if err := s.db.Preload("Package").First(&user, user.ID).Error; err != nil {
		s.logger.Error("Failed to load package info: " + err.Error())
	}

	s.logger.Info(fmt.Sprintf("Account updated successfully for user: %s", user.Username))
	return &user, nil
}

func (s *AccountService) DeleteAccount(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return errors.New("account not found")
	}

	// Begin transaction
	tx := s.db.Begin()

	// Delete related records
	if err := tx.Where("user_id = ?", id).Delete(&models.Domain{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete domains: %v", err)
	}

	if err := tx.Where("user_id = ?", id).Delete(&models.Database{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete databases: %v", err)
	}

	if err := tx.Where("user_id = ?", id).Delete(&models.Email{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete emails: %v", err)
	}

	// Delete user
	if err := tx.Delete(&user).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete user: %v", err)
	}

	tx.Commit()
	s.logger.Info(fmt.Sprintf("Account deleted successfully for user: %s", user.Username))
	return nil
}

func (s *AccountService) SuspendAccount(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return errors.New("account not found")
	}

	user.Status = "suspended"
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to suspend account: %v", err)
	}

	s.logger.Info(fmt.Sprintf("Account suspended for user: %s", user.Username))
	return nil
}

func (s *AccountService) UnsuspendAccount(id uint) error {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		return errors.New("account not found")
	}

	user.Status = "active"
	if err := s.db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to unsuspend account: %v", err)
	}

	s.logger.Info(fmt.Sprintf("Account unsuspended for user: %s", user.Username))
	return nil
}

func (s *AccountService) GetAccountByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Package").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AccountService) GetAccountByUsername(username string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).Preload("Package").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
