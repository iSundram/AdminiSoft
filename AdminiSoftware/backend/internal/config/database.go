
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
		&models.EmailAccount{},
		&models.SSLCertificate{},
		&models.Package{},
		&models.Backup{},
		&models.DNSZone{},
		&models.UserStats{},
		&models.SystemStat{},
	)
	
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
