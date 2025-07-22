
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
package models

import (
	"time"
	"gorm.io/gorm"
)

type Database struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	UserID     uint           `json:"user_id"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Name       string         `json:"name" gorm:"size:255"`
	Type       string         `json:"type" gorm:"size:20"` // mysql, postgresql, mongodb
	Host       string         `json:"host" gorm:"size:255;default:localhost"`
	Port       int            `json:"port" gorm:"default:3306"`
	Username   string         `json:"username" gorm:"size:255"`
	Password   string         `json:"-" gorm:"size:255"`
	Size       int64          `json:"size" gorm:"default:0"`
	MaxSize    int64          `json:"max_size" gorm:"default:0"`
	Status     string         `json:"status" gorm:"size:20;default:active"`
	Users      []DatabaseUser `json:"users,omitempty" gorm:"foreignKey:DatabaseID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

type DatabaseUser struct {
	ID         uint           `json:"id" gorm:"primarykey"`
	DatabaseID uint           `json:"database_id"`
	Database   Database       `json:"database,omitempty" gorm:"foreignKey:DatabaseID"`
	Username   string         `json:"username" gorm:"size:255"`
	Password   string         `json:"-" gorm:"size:255"`
	Host       string         `json:"host" gorm:"size:255;default:%"`
	Privileges string         `json:"privileges" gorm:"type:text"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
