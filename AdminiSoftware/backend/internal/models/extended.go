
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
