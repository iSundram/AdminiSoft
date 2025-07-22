
package config

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"AdminiSoftware/internal/models"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	
	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Domain{},
		&models.Database{},
		&models.Email{},
		&models.SSL{},
		&models.Package{},
		&models.Backup{},
		&models.DNS{},
		&models.Stats{},
		&models.System{},
	)
	
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
