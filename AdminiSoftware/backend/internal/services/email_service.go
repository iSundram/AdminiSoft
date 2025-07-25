
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
