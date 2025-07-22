
package whm

import (
	"AdminiSoftware/internal/models"
	"gorm.io/gorm"
)

type Manager struct {
	db     *gorm.DB
	client *Client
}

func NewManager(db *gorm.DB, client *Client) *Manager {
	return &Manager{
		db:     db,
		client: client,
	}
}

func (m *Manager) SyncAccounts() error {
	accounts, err := m.client.ListAccounts()
	if err != nil {
		return err
	}
	
	for _, account := range accounts {
		var user models.User
		result := m.db.Where("username = ?", account.Username).First(&user)
		
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			user = models.User{
				Username: account.Username,
				Email:    account.Email,
				Role:     "user",
				Status:   getStatus(account.Suspended),
			}
			m.db.Create(&user)
		} else {
			// Update existing user
			user.Status = getStatus(account.Suspended)
			m.db.Save(&user)
		}
	}
	
	return nil
}

func (m *Manager) CreateAccountFromUser(user *models.User, plan string) error {
	params := CreateAccountParams{
		Username: user.Username,
		Domain:   user.Username + ".example.com",
		Plan:     plan,
		Password: "temp_password", // Should be generated securely
	}
	
	return m.client.CreateAccount(params)
}

func (m *Manager) SuspendUser(userID uint) error {
	var user models.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return err
	}
	
	if err := m.client.SuspendAccount(user.Username); err != nil {
		return err
	}
	
	user.Status = "suspended"
	return m.db.Save(&user).Error
}

func (m *Manager) UnsuspendUser(userID uint) error {
	var user models.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return err
	}
	
	if err := m.client.UnsuspendAccount(user.Username); err != nil {
		return err
	}
	
	user.Status = "active"
	return m.db.Save(&user).Error
}

func (m *Manager) TerminateUser(userID uint) error {
	var user models.User
	if err := m.db.First(&user, userID).Error; err != nil {
		return err
	}
	
	if err := m.client.TerminateAccount(user.Username); err != nil {
		return err
	}
	
	return m.db.Delete(&user).Error
}

func getStatus(suspended bool) string {
	if suspended {
		return "suspended"
	}
	return "active"
}
