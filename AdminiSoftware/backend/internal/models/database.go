
package models

import (
	"gorm.io/gorm"
	"time"
)

type Database struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Type      string         `json:"type"`
	UserID    uint           `json:"user_id"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Size      int64          `json:"size"`
	Status    string         `json:"status"`
	Host      string         `json:"host"`
	Port      int            `json:"port"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type DatabaseUser struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Username   string         `json:"username" gorm:"not null"`
	DatabaseID uint           `json:"database_id"`
	Database   Database       `json:"database" gorm:"foreignKey:DatabaseID"`
	Privileges string         `json:"privileges"`
	Host       string         `json:"host"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
