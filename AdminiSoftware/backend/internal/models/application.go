
package models

import (
	"time"
	"gorm.io/gorm"
)

type Application struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	DomainID    uint      `json:"domain_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Path        string    `json:"path" gorm:"default:'/'"`
	DatabaseID  uint      `json:"database_id"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	Version     string    `json:"version"`
	AdminUser   string    `json:"admin_user"`
	AdminEmail  string    `json:"admin_email"`
	SiteTitle   string    `json:"site_title"`
	Description string    `json:"description"`
	Settings    string    `json:"settings" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	User     User     `json:"user" gorm:"foreignKey:UserID"`
	Domain   Domain   `json:"domain" gorm:"foreignKey:DomainID"`
	Database Database `json:"database" gorm:"foreignKey:DatabaseID"`
}

type WordPressSite struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	DomainID    uint      `json:"domain_id" gorm:"not null"`
	Path        string    `json:"path" gorm:"default:'/'"`
	SiteTitle   string    `json:"site_title" gorm:"not null"`
	AdminUser   string    `json:"admin_user" gorm:"not null"`
	AdminEmail  string    `json:"admin_email" gorm:"not null"`
	DatabaseID  uint      `json:"database_id" gorm:"not null"`
	Version     string    `json:"version" gorm:"default:'latest'"`
	Theme       string    `json:"theme" gorm:"default:'twentytwentyfour'"`
	Status      string    `json:"status" gorm:"default:'active'"`
	SiteURL     string    `json:"site_url"`
	AdminURL    string    `json:"admin_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	User     User              `json:"user" gorm:"foreignKey:UserID"`
	Domain   Domain            `json:"domain" gorm:"foreignKey:DomainID"`
	Database Database          `json:"database" gorm:"foreignKey:DatabaseID"`
	Plugins  []WordPressPlugin `json:"plugins" gorm:"foreignKey:SiteID"`
	Themes   []WordPressTheme  `json:"themes" gorm:"foreignKey:SiteID"`
}

type WordPressPlugin struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	SiteID      uint      `json:"site_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Version     string    `json:"version"`
	Status      string    `json:"status" gorm:"default:'inactive'"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Settings    string    `json:"settings" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Site WordPressSite `json:"site" gorm:"foreignKey:SiteID"`
}

type WordPressTheme struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	SiteID      uint      `json:"site_id" gorm:"not null"`
	Name        string    `json:"name" gorm:"not null"`
	Version     string    `json:"version"`
	Status      string    `json:"status" gorm:"default:'inactive'"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Screenshot  string    `json:"screenshot"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// Relationships
	Site WordPressSite `json:"site" gorm:"foreignKey:SiteID"`
}

func (Application) TableName() string {
	return "applications"
}

func (WordPressSite) TableName() string {
	return "wordpress_sites"
}

func (WordPressPlugin) TableName() string {
	return "wordpress_plugins"
}

func (WordPressTheme) TableName() string {
	return "wordpress_themes"
}

// Hooks
func (a *Application) BeforeCreate(tx *gorm.DB) (err error) {
	if a.Path == "" {
		a.Path = "/"
	}
	return
}

func (w *WordPressSite) BeforeCreate(tx *gorm.DB) (err error) {
	if w.Path == "" {
		w.Path = "/"
	}
	if w.Version == "" {
		w.Version = "latest"
	}
	if w.Theme == "" {
		w.Theme = "twentytwentyfour"
	}
	return
}
package models

import (
	"time"
	"gorm.io/gorm"
)

type Application struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	UserID      uint           `json:"user_id"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DomainID    uint           `json:"domain_id"`
	Domain      Domain         `json:"domain,omitempty" gorm:"foreignKey:DomainID"`
	Name        string         `json:"name" gorm:"size:255"`
	Type        string         `json:"type" gorm:"size:50"` // wordpress, drupal, nodejs, python, etc.
	Version     string         `json:"version" gorm:"size:50"`
	Path        string         `json:"path" gorm:"size:500"`
	URL         string         `json:"url" gorm:"size:500"`
	Status      string         `json:"status" gorm:"size:20;default:active"`
	Config      string         `json:"config" gorm:"type:text"`
	Variables   []AppVariable  `json:"variables,omitempty" gorm:"foreignKey:ApplicationID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type AppVariable struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	ApplicationID uint           `json:"application_id"`
	Application   Application    `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	Key           string         `json:"key" gorm:"size:255"`
	Value         string         `json:"value" gorm:"type:text"`
	IsSecret      bool           `json:"is_secret" gorm:"default:false"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

type WordPressInstall struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	ApplicationID uint           `json:"application_id"`
	Application   Application    `json:"application,omitempty" gorm:"foreignKey:ApplicationID"`
	AdminUser     string         `json:"admin_user" gorm:"size:255"`
	AdminEmail    string         `json:"admin_email" gorm:"size:255"`
	SiteTitle    string         `json:"site_title" gorm:"size:255"`
	Plugins       string         `json:"plugins" gorm:"type:text"`
	Themes        string         `json:"themes" gorm:"type:text"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
