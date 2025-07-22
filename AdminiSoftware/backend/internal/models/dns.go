
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
