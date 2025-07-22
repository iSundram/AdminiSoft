
package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Username    string    `json:"username" gorm:"uniqueIndex;not null"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null"`
	Password    string    `json:"-" gorm:"not null"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Role        string    `json:"role" gorm:"default:'user'"` // admin, reseller, user
	Status      string    `json:"status" gorm:"default:'active'"` // active, suspended, deleted
	PackageID   *uint     `json:"package_id"`
	Package     *Package  `json:"package,omitempty"`
	ResellerID  *uint     `json:"reseller_id"`
	Reseller    *User     `json:"reseller,omitempty"`
	DiskUsed    int64     `json:"disk_used" gorm:"default:0"`
	BandwidthUsed int64   `json:"bandwidth_used" gorm:"default:0"`
	LastLogin   *time.Time `json:"last_login"`
	TwoFactorEnabled bool `json:"two_factor_enabled" gorm:"default:false"`
	TwoFactorSecret string `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Relationships
	Domains     []Domain   `json:"domains,omitempty"`
	Databases   []Database `json:"databases,omitempty"`
	Emails      []Email    `json:"emails,omitempty"`
}

type LoginAttempt struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	IP        string    `json:"ip" gorm:"index"`
	Username  string    `json:"username" gorm:"index"`
	Success   bool      `json:"success"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}
