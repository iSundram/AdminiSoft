
package models

import (
	"time"
	"gorm.io/gorm"
)

type EmailForwarder struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	DomainID    uint      `json:"domain_id" gorm:"not null"`
	Source      string    `json:"source" gorm:"not null"`
	Destination string    `json:"destination" gorm:"not null"`
	Status      string    `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}

type BandwidthStat struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	DomainID  uint      `json:"domain_id"`
	Date      time.Time `json:"date" gorm:"not null"`
	Incoming  int64     `json:"incoming"`
	Outgoing  int64     `json:"outgoing"`
	Total     int64     `json:"total"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}

type VisitorStat struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	DomainID     uint      `json:"domain_id"`
	Date         time.Time `json:"date" gorm:"not null"`
	Visitors     int       `json:"visitors"`
	PageViews    int       `json:"page_views"`
	UniqueVisits int       `json:"unique_visits"`
	Bandwidth    int64     `json:"bandwidth"`
	CreatedAt    time.Time `json:"created_at"`
	
	// Relationships
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}

type ErrorLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	DomainID  uint      `json:"domain_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message" gorm:"type:text"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}

type AccessLog struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	DomainID   uint      `json:"domain_id"`
	IP         string    `json:"ip"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	StatusCode int       `json:"status_code"`
	Size       int64     `json:"size"`
	Referer    string    `json:"referer"`
	UserAgent  string    `json:"user_agent"`
	Timestamp  time.Time `json:"timestamp"`
	
	// Relationships
	User   User   `json:"user" gorm:"foreignKey:UserID"`
	Domain Domain `json:"domain" gorm:"foreignKey:DomainID"`
}

type ResourceUsage struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CPUUsage  float64   `json:"cpu_usage"`
	MemUsage  int64     `json:"memory_usage"`
	DiskIO    int64     `json:"disk_io"`
	NetworkIO int64     `json:"network_io"`
	Processes int       `json:"processes"`
	Timestamp time.Time `json:"timestamp"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type SecurityEvent struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id"`
	Type        string    `json:"type" gorm:"not null"`
	Description string    `json:"description" gorm:"type:text"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Severity    string    `json:"severity" gorm:"default:'info'"`
	Status      string    `json:"status" gorm:"default:'active'"`
	Data        string    `json:"data" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Action      string    `json:"action" gorm:"not null"`
	Resource    string    `json:"resource"`
	ResourceID  uint      `json:"resource_id"`
	Description string    `json:"description"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	Data        string    `json:"data" gorm:"type:json"`
	CreatedAt   time.Time `json:"created_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

type CronJob struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	Command     string    `json:"command" gorm:"not null"`
	Schedule    string    `json:"schedule" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'active'"`
	LastRun     *time.Time `json:"last_run"`
	NextRun     *time.Time `json:"next_run"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}

// Table names
func (EmailForwarder) TableName() string {
	return "email_forwarders"
}

func (BandwidthStat) TableName() string {
	return "bandwidth_stats"
}

func (VisitorStat) TableName() string {
	return "visitor_stats"
}

func (ErrorLog) TableName() string {
	return "error_logs"
}

func (AccessLog) TableName() string {
	return "access_logs"
}

func (ResourceUsage) TableName() string {
	return "resource_usage"
}

func (SecurityEvent) TableName() string {
	return "security_events"
}

func (ActivityLog) TableName() string {
	return "activity_logs"
}

func (CronJob) TableName() string {
	return "cron_jobs"
}
package models

import "time"

// Account Management Requests
type CreateAccountRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=32"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	PackageID uint   `json:"package_id" binding:"required"`
	Domain    string `json:"domain" binding:"required"`
}

type UpdateAccountRequest struct {
	Email     string `json:"email,omitempty" binding:"omitempty,email"`
	PackageID uint   `json:"package_id,omitempty"`
	Status    string `json:"status,omitempty"`
}

// Package Management Requests
type CreatePackageRequest struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	DiskQuota      int64   `json:"disk_quota" binding:"required"`
	BandwidthQuota int64   `json:"bandwidth_quota" binding:"required"`
	EmailAccounts  int     `json:"email_accounts" binding:"required"`
	Databases      int     `json:"databases" binding:"required"`
	Subdomains     int     `json:"subdomains" binding:"required"`
	FTPAccounts    int     `json:"ftp_accounts" binding:"required"`
	Price          float64 `json:"price" binding:"required"`
	BillingCycle   string  `json:"billing_cycle" binding:"required"`
	Features       string  `json:"features"`
}

type UpdatePackageRequest struct {
	Name           string  `json:"name,omitempty"`
	Description    string  `json:"description,omitempty"`
	DiskQuota      int64   `json:"disk_quota,omitempty"`
	BandwidthQuota int64   `json:"bandwidth_quota,omitempty"`
	EmailAccounts  int     `json:"email_accounts,omitempty"`
	Databases      int     `json:"databases,omitempty"`
	Subdomains     int     `json:"subdomains,omitempty"`
	FTPAccounts    int     `json:"ftp_accounts,omitempty"`
	Price          float64 `json:"price,omitempty"`
	BillingCycle   string  `json:"billing_cycle,omitempty"`
	Features       string  `json:"features,omitempty"`
}

// Domain Management Requests
type CreateDomainRequest struct {
	Domain     string `json:"domain" binding:"required"`
	DocumentRoot string `json:"document_root"`
	Status     string `json:"status"`
}

type UpdateDomainRequest struct {
	DocumentRoot string `json:"document_root,omitempty"`
	Status       string `json:"status,omitempty"`
}

// Database Management Requests
type CreateDatabaseRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateDatabaseRequest struct {
	Password string `json:"password,omitempty"`
	Status   string `json:"status,omitempty"`
}

// Email Management Requests
type CreateEmailRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Quota    int64  `json:"quota"`
}

type UpdateEmailRequest struct {
	Password string `json:"password,omitempty"`
	Quota    int64  `json:"quota,omitempty"`
	Status   string `json:"status,omitempty"`
}

// SSL Certificate Requests
type CreateSSLRequest struct {
	Domain       string `json:"domain" binding:"required"`
	Certificate  string `json:"certificate" binding:"required"`
	PrivateKey   string `json:"private_key" binding:"required"`
	Intermediate string `json:"intermediate"`
}

// Backup Management Requests
type CreateBackupRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
}

// DNS Management Requests
type CreateDNSRecordRequest struct {
	Type     string `json:"type" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Content  string `json:"content" binding:"required"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
}

type UpdateDNSRecordRequest struct {
	Content  string `json:"content,omitempty"`
	TTL      int    `json:"ttl,omitempty"`
	Priority int    `json:"priority,omitempty"`
}

// System Settings
type SystemSetting struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Key         string    `json:"key" gorm:"uniqueIndex;not null"`
	Value       string    `json:"value" gorm:"not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Activity Log
type ActivityLog struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Action      string    `json:"action" gorm:"not null"`
	Resource    string    `json:"resource"`
	ResourceID  string    `json:"resource_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// API Response structures
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

// Dashboard Statistics
type DashboardStats struct {
	TotalUsers     int64   `json:"total_users"`
	ActiveUsers    int64   `json:"active_users"`
	TotalDomains   int64   `json:"total_domains"`
	TotalDatabases int64   `json:"total_databases"`
	TotalEmails    int64   `json:"total_emails"`
	DiskUsage      int64   `json:"disk_usage"`
	BandwidthUsage int64   `json:"bandwidth_usage"`
	ServerLoad     float64 `json:"server_load"`
	Uptime         string  `json:"uptime"`
}

// Resource Usage
type ResourceUsage struct {
	UserID         uint    `json:"user_id"`
	DiskUsed       int64   `json:"disk_used"`
	DiskLimit      int64   `json:"disk_limit"`
	BandwidthUsed  int64   `json:"bandwidth_used"`
	BandwidthLimit int64   `json:"bandwidth_limit"`
	EmailsUsed     int     `json:"emails_used"`
	EmailsLimit    int     `json:"emails_limit"`
	DatabasesUsed  int     `json:"databases_used"`
	DatabasesLimit int     `json:"databases_limit"`
	LastUpdated    time.Time `json:"last_updated"`
}
