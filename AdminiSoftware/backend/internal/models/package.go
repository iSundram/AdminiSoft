
package models

import (
	"time"
	"gorm.io/gorm"
)

type Package struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"uniqueIndex;not null"`
	Description     string    `json:"description"`
	DiskQuota       int64     `json:"disk_quota"` // MB, 0 = unlimited
	BandwidthQuota  int64     `json:"bandwidth_quota"` // MB, 0 = unlimited
	EmailAccounts   int       `json:"email_accounts"` // 0 = unlimited
	Databases       int       `json:"databases"` // 0 = unlimited
	Subdomains      int       `json:"subdomains"` // 0 = unlimited
	ParkedDomains   int       `json:"parked_domains"` // 0 = unlimited
	AddonDomains    int       `json:"addon_domains"` // 0 = unlimited
	FTPAccounts     int       `json:"ftp_accounts"` // 0 = unlimited
	SSLCertificates int       `json:"ssl_certificates"` // 0 = unlimited
	CGIAccess       bool      `json:"cgi_access" gorm:"default:true"`
	SSHAccess       bool      `json:"ssh_access" gorm:"default:false"`
	CronJobs        int       `json:"cron_jobs"` // 0 = unlimited
	Price           float64   `json:"price" gorm:"default:0"`
	Currency        string    `json:"currency" gorm:"default:'USD'"`
	BillingCycle    string    `json:"billing_cycle" gorm:"default:'monthly'"` // monthly, yearly
	Status          string    `json:"status" gorm:"default:'active'"` // active, inactive
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
	
	// Features
	Features        []PackageFeature `json:"features,omitempty"`
}

type PackageFeature struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PackageID uint      `json:"package_id"`
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled" gorm:"default:true"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
package models

import (
	"time"
	"gorm.io/gorm"
)

type Package struct {
	ID                uint           `json:"id" gorm:"primarykey"`
	Name              string         `json:"name" gorm:"uniqueIndex;size:255"`
	DiskQuotaMB       int            `json:"disk_quota_mb" gorm:"default:1000"`
	BandwidthMB       int            `json:"bandwidth_mb" gorm:"default:10000"`
	EmailAccounts     int            `json:"email_accounts" gorm:"default:10"`
	Databases         int            `json:"databases" gorm:"default:5"`
	SubDomains        int            `json:"sub_domains" gorm:"default:10"`
	ParkedDomains     int            `json:"parked_domains" gorm:"default:5"`
	AddonDomains      int            `json:"addon_domains" gorm:"default:5"`
	FTPAccounts       int            `json:"ftp_accounts" gorm:"default:5"`
	CronJobs          int            `json:"cron_jobs" gorm:"default:5"`
	CGIAccess         bool           `json:"cgi_access" gorm:"default:false"`
	SSHAccess         bool           `json:"ssh_access" gorm:"default:false"`
	SSLSupport        bool           `json:"ssl_support" gorm:"default:true"`
	Features          string         `json:"features" gorm:"type:text"`
	Status            string         `json:"status" gorm:"size:20;default:active"`
	Users             []User         `json:"users,omitempty" gorm:"foreignKey:PackageID"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`
}
