
package models

import (
	"gorm.io/gorm"
	"time"
)

type Backup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Type        string         `json:"type"`
	Status      string         `json:"status"`
	Size        int64          `json:"size"`
	Path        string         `json:"path"`
	RemotePath  string         `json:"remote_path"`
	Compressed  bool           `json:"compressed"`
	Encrypted   bool           `json:"encrypted"`
	Description string         `json:"description"`
	StartedAt   *time.Time     `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type BackupSchedule struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Name        string         `json:"name" gorm:"not null"`
	Type        string         `json:"type"`
	Frequency   string         `json:"frequency"`
	Time        string         `json:"time"`
	Retention   int            `json:"retention"`
	Enabled     bool           `json:"enabled"`
	LastRun     *time.Time     `json:"last_run"`
	NextRun     *time.Time     `json:"next_run"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
package models

import (
	"time"
	"gorm.io/gorm"
)

type Backup struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type        string         `json:"type" gorm:"size:20"` // full, partial, database, files
	Status      string         `json:"status" gorm:"size:20;default:pending"`
	Path        string         `json:"path" gorm:"size:500"`
	SizeMB      int64          `json:"size_mb" gorm:"default:0"`
	Description string         `json:"description" gorm:"type:text"`
	Progress    int            `json:"progress" gorm:"default:0"`
	ErrorLog    string         `json:"error_log" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type BackupSchedule struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Type      string         `json:"type" gorm:"size:20"`
	Schedule  string         `json:"schedule" gorm:"size:100"` // cron format
	Retention int            `json:"retention" gorm:"default:7"` // days
	Status    string         `json:"status" gorm:"size:20;default:active"`
	LastRun   *time.Time     `json:"last_run"`
	NextRun   *time.Time     `json:"next_run"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
