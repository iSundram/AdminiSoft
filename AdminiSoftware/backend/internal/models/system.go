
package models

import (
	"gorm.io/gorm"
	"time"
)

type SystemStat struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CPUUsage  float64        `json:"cpu_usage"`
	RAMUsage  float64        `json:"ram_usage"`
	DiskUsage float64        `json:"disk_usage"`
	LoadAvg   string         `json:"load_avg"`
	Timestamp time.Time      `json:"timestamp"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type SystemService struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Status      string         `json:"status"`
	Port        int            `json:"port"`
	Description string         `json:"description"`
	AutoStart   bool           `json:"auto_start"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type ServerConfig struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Key         string         `json:"key" gorm:"uniqueIndex;not null"`
	Value       string         `json:"value"`
	Type        string         `json:"type"`
	Category    string         `json:"category"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Branding struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user" gorm:"foreignKey:UserID"`
	LogoURL      string         `json:"logo_url"`
	CompanyName  string         `json:"company_name"`
	SupportURL   string         `json:"support_url"`
	TermsURL     string         `json:"terms_url"`
	ThemeColor   string         `json:"theme_color"`
	CustomCSS    string         `json:"custom_css"`
	ShowPoweredBy bool          `json:"show_powered_by"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
