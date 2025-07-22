
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
package config

import (
	"AdminiSoftware/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDatabase(cfg *Config) (*gorm.DB, error) {
	var logLevel logger.LogLevel = logger.Silent
	if cfg.Debug {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
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
		&models.Application{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
