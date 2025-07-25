
package models

import (
	"gorm.io/gorm"
	"time"
)

type SSLCertificate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Domain      string         `json:"domain" gorm:"not null"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Type        string         `json:"type"`
	Provider    string         `json:"provider"`
	Certificate string         `json:"certificate"`
	PrivateKey  string         `json:"private_key"`
	IssuedAt    *time.Time     `json:"issued_at"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	AutoRenew   bool           `json:"auto_renew"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
