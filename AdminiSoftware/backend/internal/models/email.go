
package models

import (
	"gorm.io/gorm"
	"time"
)

type EmailAccount struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Password  string         `json:"password"`
	Quota     int64          `json:"quota"`
	Used      int64          `json:"used"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type EmailForwarder struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Source      string         `json:"source" gorm:"not null"`
	Destination string         `json:"destination" gorm:"not null"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type MailingList struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Members     []string       `json:"members" gorm:"type:text[]"`
	Status      string         `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
package models

import (
	"time"
	"gorm.io/gorm"
)

type Email struct {
	ID           uint           `json:"id" gorm:"primarykey"`
	UserID       uint           `json:"user_id"`
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DomainID     uint           `json:"domain_id"`
	Domain       Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Email        string         `json:"email" gorm:"uniqueIndex;size:255"`
	Password     string         `json:"-" gorm:"size:255"`
	QuotaMB      int            `json:"quota_mb" gorm:"default:100"`
	UsedMB       int            `json:"used_mb" gorm:"default:0"`
	Status       string         `json:"status" gorm:"size:20;default:active"`
	Forwarders   []EmailForwarder `json:"forwarders,omitempty" gorm:"foreignKey:EmailID"`
	Filters      []EmailFilter  `json:"filters,omitempty" gorm:"foreignKey:EmailID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type EmailForwarder struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	EmailID     uint           `json:"email_id"`
	Email       Email          `json:"email,omitempty" gorm:"foreignKey:EmailID"`
	Source      string         `json:"source" gorm:"size:255"`
	Destination string         `json:"destination" gorm:"size:255"`
	Status      string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type EmailFilter struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	EmailID   uint           `json:"email_id"`
	Email     Email          `json:"email,omitempty" gorm:"foreignKey:EmailID"`
	Name      string         `json:"name" gorm:"size:255"`
	Condition string         `json:"condition" gorm:"type:text"`
	Action    string         `json:"action" gorm:"type:text"`
	Status    string         `json:"status" gorm:"size:20;default:active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
