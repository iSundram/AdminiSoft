
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
package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	Username          string         `json:"username" gorm:"uniqueIndex;size:255"`
	Email            string         `json:"email" gorm:"uniqueIndex;size:255"`
	Password         string         `json:"-" gorm:"size:255"`
	Role             string         `json:"role" gorm:"size:20;default:user"`
	Status           string         `json:"status" gorm:"size:20;default:active"`
	FirstName        string         `json:"first_name" gorm:"size:100"`
	LastName         string         `json:"last_name" gorm:"size:100"`
	ContactEmail     string         `json:"contact_email" gorm:"size:255"`
	Theme            string         `json:"theme" gorm:"size:50;default:default"`
	Language         string         `json:"language" gorm:"size:10;default:en"`
	TwoFactorEnabled bool           `json:"two_factor_enabled" gorm:"default:false"`
	TwoFactorSecret  string         `json:"-" gorm:"size:255"`
	LastLogin        *time.Time     `json:"last_login"`
	PackageID        *uint          `json:"package_id"`
	Package          *Package       `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	Domains          []Domain       `json:"domains,omitempty" gorm:"foreignKey:UserID"`
	Databases        []Database     `json:"databases,omitempty" gorm:"foreignKey:UserID"`
	Emails           []Email        `json:"emails,omitempty" gorm:"foreignKey:UserID"`
	SSLs             []SSL          `json:"ssls,omitempty" gorm:"foreignKey:UserID"`
	Backups          []Backup       `json:"backups,omitempty" gorm:"foreignKey:UserID"`
	Applications     []Application  `json:"applications,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

type LoginAttempt struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"size:255"`
	IP        string    `json:"ip" gorm:"size:45"`
	Success   bool      `json:"success"`
	UserAgent string    `json:"user_agent" gorm:"size:500"`
	CreatedAt time.Time `json:"created_at"`
}
