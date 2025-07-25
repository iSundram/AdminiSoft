
package models

import (
	"gorm.io/gorm"
	"time"
)

type UserStats struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserID          uint           `json:"user_id"`
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	BandwidthUsed   int64          `json:"bandwidth_used"`
	DiskUsed        int64          `json:"disk_used"`
	EmailsSent      int            `json:"emails_sent"`
	DomainsCount    int            `json:"domains_count"`
	DatabasesCount  int            `json:"databases_count"`
	EmailsCount     int            `json:"emails_count"`
	Date            time.Time      `json:"date"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type AccessLog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	IP        string         `json:"ip"`
	UserAgent string         `json:"user_agent"`
	Action    string         `json:"action"`
	Resource  string         `json:"resource"`
	Status    string         `json:"status"`
	Timestamp time.Time      `json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
