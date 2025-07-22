
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
package models

import (
	"time"
	"gorm.io/gorm"
)

type Stats struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        uint           `json:"user_id"`
	User          User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Date          time.Time      `json:"date"`
	DiskUsageMB   int64          `json:"disk_usage_mb" gorm:"default:0"`
	BandwidthMB   int64          `json:"bandwidth_mb" gorm:"default:0"`
	CPUUsage      float64        `json:"cpu_usage" gorm:"default:0"`
	MemoryUsageMB int64          `json:"memory_usage_mb" gorm:"default:0"`
	Visits        int64          `json:"visits" gorm:"default:0"`
	PageViews     int64          `json:"page_views" gorm:"default:0"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type BandwidthUsage struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Date      time.Time `json:"date"`
	HTTP      int64     `json:"http" gorm:"default:0"`
	HTTPS     int64     `json:"https" gorm:"default:0"`
	FTP       int64     `json:"ftp" gorm:"default:0"`
	Mail      int64     `json:"mail" gorm:"default:0"`
	Total     int64     `json:"total" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
