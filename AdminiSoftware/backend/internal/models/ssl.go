
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
package models

import (
	"time"
	"gorm.io/gorm"
)

type SSL struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DomainID    uint           `json:"domain_id"`
	Domain      Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Type        string         `json:"type" gorm:"size:50"` // letsencrypt, self-signed, commercial
	Certificate string         `json:"certificate" gorm:"type:text"`
	PrivateKey  string         `json:"private_key" gorm:"type:text"`
	Chain       string         `json:"chain" gorm:"type:text"`
	Status      string         `json:"status" gorm:"size:20;default:active"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	AutoRenew   bool           `json:"auto_renew" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
