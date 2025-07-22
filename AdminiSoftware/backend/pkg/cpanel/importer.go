
package cpanel

import (
	"AdminiSoftware/internal/models"
	"encoding/json"
	"gorm.io/gorm"
)

type Importer struct {
	db     *gorm.DB
	client *Client
}

func NewImporter(db *gorm.DB, client *Client) *Importer {
	return &Importer{
		db:     db,
		client: client,
	}
}

func (i *Importer) ImportAccount(cpanelUser string) error {
	// Import user account from cPanel
	user := &models.User{
		Username: cpanelUser,
		Email:    cpanelUser + "@example.com",
		Role:     "user",
		Status:   "active",
	}
	
	if err := i.db.Create(user).Error; err != nil {
		return err
	}
	
	// Import domains
	domains, err := i.client.ListDomains()
	if err != nil {
		return err
	}
	
	for _, domainName := range domains {
		domain := &models.Domain{
			UserID:      user.ID,
			Name:        domainName,
			Type:        "primary",
			Status:      "active",
			DocumentRoot: "/public_html",
		}
		i.db.Create(domain)
	}
	
	return nil
}

func (i *Importer) ImportEmailAccounts(userID uint) error {
	// Import email accounts from cPanel
	return nil
}

func (i *Importer) ImportDatabases(userID uint) error {
	// Import databases from cPanel
	return nil
}

func (i *Importer) ExportAccount(userID uint) (map[string]interface{}, error) {
	var user models.User
	if err := i.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	
	var domains []models.Domain
	i.db.Where("user_id = ?", userID).Find(&domains)
	
	export := map[string]interface{}{
		"user":    user,
		"domains": domains,
	}
	
	return export, nil
}

func (i *Importer) ExportToJSON(userID uint) ([]byte, error) {
	data, err := i.ExportAccount(userID)
	if err != nil {
		return nil, err
	}
	
	return json.MarshalIndent(data, "", "  ")
}
