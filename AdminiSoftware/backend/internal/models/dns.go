
package models

import (
	"gorm.io/gorm"
	"time"
)

type DNSZone struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type"`
	TTL       int            `json:"ttl"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Records   []DNSRecord    `json:"records" gorm:"foreignKey:ZoneID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type DNSRecord struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ZoneID    uint           `json:"zone_id"`
	Zone      DNSZone        `json:"zone" gorm:"foreignKey:ZoneID"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type" gorm:"not null"`
	Value     string         `json:"value" gorm:"not null"`
	TTL       int            `json:"ttl"`
	Priority  int            `json:"priority"`
	Weight    int            `json:"weight"`
	Port      int            `json:"port"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
package models

import (
	"time"
	"gorm.io/gorm"
)

type DNS struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	DomainID  uint           `json:"domain_id"`
	Domain    Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Name      string         `json:"name" gorm:"size:255"`
	Type      string         `json:"type" gorm:"size:10"` // A, AAAA, CNAME, MX, TXT, etc.
	Value     string         `json:"value" gorm:"size:500"`
	TTL       int            `json:"ttl" gorm:"default:14400"`
	Priority  int            `json:"priority" gorm:"default:0"`
	Status    string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DNSZone struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	DomainID  uint           `json:"domain_id"`
	Domain    Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Serial    string         `json:"serial" gorm:"size:20"`
	Refresh   int            `json:"refresh" gorm:"default:86400"`
	Retry     int            `json:"retry" gorm:"default:7200"`
	Expire    int            `json:"expire" gorm:"default:3600000"`
	Minimum   int            `json:"minimum" gorm:"default:86400"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
