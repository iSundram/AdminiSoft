
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
package models

import (
	"time"
	"gorm.io/gorm"
)

type Domain struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name         string         `json:"name" gorm:"uniqueIndex;size:255"`
	Type         string         `json:"type" gorm:"size:20"` // primary, addon, subdomain, parked
	Status       string         `json:"status" gorm:"size:20;default:active"`
	DocumentRoot string         `json:"document_root" gorm:"size:500"`
	Redirects    []DomainRedirect `json:"redirects,omitempty" gorm:"foreignKey:DomainID"`
	DNSRecords   []DNS          `json:"dns_records,omitempty" gorm:"foreignKey:DomainID"`
	SSLEnabled   bool           `json:"ssl_enabled" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type DomainRedirect struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	DomainID  uint           `json:"domain_id"`
	Domain    Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Source    string         `json:"source" gorm:"size:255"`
	Target    string         `json:"target" gorm:"size:500"`
	Type      string         `json:"type" gorm:"size:20;default:301"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
