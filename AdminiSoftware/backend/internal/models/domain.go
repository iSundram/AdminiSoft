
package models

import (
	"time"
	"gorm.io/gorm"
)

type Domain struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	User        User      `json:"user"`
	Type        string    `json:"type" gorm:"default:'addon'"` // main, addon, subdomain, parked
	DocumentRoot string   `json:"document_root"`
	Status      string    `json:"status" gorm:"default:'active'"` // active, suspended, deleted
	SSLEnabled  bool      `json:"ssl_enabled" gorm:"default:false"`
	AutoSSL     bool      `json:"auto_ssl" gorm:"default:true"`
	BandwidthLimit int64  `json:"bandwidth_limit" gorm:"default:0"` // 0 = unlimited
	BandwidthUsed  int64  `json:"bandwidth_used" gorm:"default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	DNS         []DNS     `json:"dns,omitempty"`
	SSL         []SSL     `json:"ssl,omitempty"`
}

type Subdomain struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	DomainID  uint      `json:"domain_id" gorm:"not null"`
	Domain    Domain    `json:"domain"`
	DocumentRoot string `json:"document_root"`
	Status    string    `json:"status" gorm:"default:'active'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
