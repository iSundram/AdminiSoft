
package services

import (
	"AdminiSoftware/internal/models"
	"fmt"
	"gorm.io/gorm"
)

type DatabaseService struct {
	db *gorm.DB
}

func NewDatabaseService(db *gorm.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

func (s *DatabaseService) CreateDatabase(database *models.Database) error {
	if err := s.db.Create(database).Error; err != nil {
		return err
	}
	
	// Create actual database
	switch database.Type {
	case "mysql":
		return s.createMySQLDatabase(database)
	case "postgresql":
		return s.createPostgreSQLDatabase(database)
	case "mongodb":
		return s.createMongoDatabase(database)
	default:
		return fmt.Errorf("unsupported database type: %s", database.Type)
	}
}

func (s *DatabaseService) GetDatabases(userID uint) ([]models.Database, error) {
	var databases []models.Database
	err := s.db.Where("user_id = ?", userID).Find(&databases).Error
	return databases, err
}

func (s *DatabaseService) DeleteDatabase(id uint) error {
	var database models.Database
	if err := s.db.First(&database, id).Error; err != nil {
		return err
	}
	
	// Delete actual database
	switch database.Type {
	case "mysql":
		s.dropMySQLDatabase(&database)
	case "postgresql":
		s.dropPostgreSQLDatabase(&database)
	case "mongodb":
		s.dropMongoDatabase(&database)
	}
	
	return s.db.Delete(&database).Error
}

func (s *DatabaseService) CreateUser(user *models.DatabaseUser) error {
	return s.db.Create(user).Error
}

func (s *DatabaseService) GetUsers(databaseID uint) ([]models.DatabaseUser, error) {
	var users []models.DatabaseUser
	err := s.db.Where("database_id = ?", databaseID).Find(&users).Error
	return users, err
}

func (s *DatabaseService) createMySQLDatabase(database *models.Database) error {
	// Implementation for creating MySQL database
	return nil
}

func (s *DatabaseService) createPostgreSQLDatabase(database *models.Database) error {
	// Implementation for creating PostgreSQL database
	return nil
}

func (s *DatabaseService) createMongoDatabase(database *models.Database) error {
	// Implementation for creating MongoDB database
	return nil
}

func (s *DatabaseService) dropMySQLDatabase(database *models.Database) error {
	// Implementation for dropping MySQL database
	return nil
}

func (s *DatabaseService) dropPostgreSQLDatabase(database *models.Database) error {
	// Implementation for dropping PostgreSQL database
	return nil
}

func (s *DatabaseService) dropMongoDatabase(database *models.Database) error {
	// Implementation for dropping MongoDB database
	return nil
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"fmt"
	"os/exec"

	"gorm.io/gorm"
)

type DatabaseService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewDatabaseService(db *gorm.DB, logger *utils.Logger) *DatabaseService {
	return &DatabaseService{
		db:     db,
		logger: logger,
	}
}

func (s *DatabaseService) CreateDatabase(userID uint, req *models.CreateDatabaseRequest) (*models.Database, error) {
	// Check if database name already exists for user
	var existing models.Database
	if err := s.db.Where("user_id = ? AND name = ?", userID, req.Name).First(&existing).Error; err == nil {
		return nil, fmt.Errorf("database name already exists")
	}

	// Create database record
	database := &models.Database{
		UserID:   userID,
		Name:     fmt.Sprintf("user%d_%s", userID, req.Name),
		Type:     req.Type,
		Username: fmt.Sprintf("user%d_%s", userID, req.Username),
		Status:   "creating",
		Size:     0,
	}

	if err := s.db.Create(database).Error; err != nil {
		return nil, fmt.Errorf("failed to create database record: %v", err)
	}

	// Create actual database
	if err := s.createPhysicalDatabase(database, req.Password); err != nil {
		s.db.Model(database).Update("status", "failed")
		return nil, fmt.Errorf("failed to create physical database: %v", err)
	}

	// Update status
	s.db.Model(database).Update("status", "active")

	s.logger.Info(fmt.Sprintf("Database created: %s for user %d", database.Name, userID))
	return database, nil
}

func (s *DatabaseService) createPhysicalDatabase(db *models.Database, password string) error {
	switch db.Type {
	case "mysql":
		return s.createMySQLDatabase(db, password)
	case "postgresql":
		return s.createPostgreSQLDatabase(db, password)
	default:
		return fmt.Errorf("unsupported database type: %s", db.Type)
	}
}

func (s *DatabaseService) createMySQLDatabase(db *models.Database, password string) error {
	// Create database
	cmd := exec.Command("mysql", "-u", "root", "-e", fmt.Sprintf("CREATE DATABASE `%s`;", db.Name))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create MySQL database: %v", err)
	}

	// Create user
	cmd = exec.Command("mysql", "-u", "root", "-e", 
		fmt.Sprintf("CREATE USER '%s'@'localhost' IDENTIFIED BY '%s';", db.Username, password))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create MySQL user: %v", err)
	}

	// Grant privileges
	cmd = exec.Command("mysql", "-u", "root", "-e", 
		fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'localhost';", db.Name, db.Username))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to grant MySQL privileges: %v", err)
	}

	// Flush privileges
	cmd = exec.Command("mysql", "-u", "root", "-e", "FLUSH PRIVILEGES;")
	return cmd.Run()
}

func (s *DatabaseService) createPostgreSQLDatabase(db *models.Database, password string) error {
	// Create user
	cmd := exec.Command("sudo", "-u", "postgres", "createuser", db.Username)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create PostgreSQL user: %v", err)
	}

	// Set password
	cmd = exec.Command("sudo", "-u", "postgres", "psql", "-c", 
		fmt.Sprintf("ALTER USER %s PASSWORD '%s';", db.Username, password))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set PostgreSQL password: %v", err)
	}

	// Create database
	cmd = exec.Command("sudo", "-u", "postgres", "createdb", "-O", db.Username, db.Name)
	return cmd.Run()
}

func (s *DatabaseService) GetDatabases(userID uint) ([]models.Database, error) {
	var databases []models.Database
	if err := s.db.Where("user_id = ?", userID).Find(&databases).Error; err != nil {
		return nil, err
	}
	return databases, nil
}

func (s *DatabaseService) GetDatabase(id uint, userID uint) (*models.Database, error) {
	var database models.Database
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		return nil, err
	}
	return &database, nil
}

func (s *DatabaseService) UpdateDatabase(id uint, userID uint, req *models.UpdateDatabaseRequest) (*models.Database, error) {
	var database models.Database
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		return nil, err
	}

	// Update password if provided
	if req.Password != "" {
		if err := s.updateDatabasePassword(&database, req.Password); err != nil {
			return nil, fmt.Errorf("failed to update database password: %v", err)
		}
	}

	// Update status if provided
	if req.Status != "" {
		database.Status = req.Status
	}

	if err := s.db.Save(&database).Error; err != nil {
		return nil, fmt.Errorf("failed to update database: %v", err)
	}

	return &database, nil
}

func (s *DatabaseService) updateDatabasePassword(db *models.Database, password string) error {
	switch db.Type {
	case "mysql":
		cmd := exec.Command("mysql", "-u", "root", "-e", 
			fmt.Sprintf("ALTER USER '%s'@'localhost' IDENTIFIED BY '%s';", db.Username, password))
		return cmd.Run()
	case "postgresql":
		cmd := exec.Command("sudo", "-u", "postgres", "psql", "-c", 
			fmt.Sprintf("ALTER USER %s PASSWORD '%s';", db.Username, password))
		return cmd.Run()
	default:
		return fmt.Errorf("unsupported database type: %s", db.Type)
	}
}

func (s *DatabaseService) DeleteDatabase(id uint, userID uint) error {
	var database models.Database
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		return err
	}

	// Delete physical database
	if err := s.deletePhysicalDatabase(&database); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete physical database %s: %v", database.Name, err))
	}

	// Delete database record
	if err := s.db.Delete(&database).Error; err != nil {
		return fmt.Errorf("failed to delete database record: %v", err)
	}

	s.logger.Info(fmt.Sprintf("Database deleted: %s for user %d", database.Name, userID))
	return nil
}

func (s *DatabaseService) deletePhysicalDatabase(db *models.Database) error {
	switch db.Type {
	case "mysql":
		// Drop database
		cmd := exec.Command("mysql", "-u", "root", "-e", fmt.Sprintf("DROP DATABASE IF EXISTS `%s`;", db.Name))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to drop MySQL database: %v", err)
		}
		
		// Drop user
		cmd = exec.Command("mysql", "-u", "root", "-e", 
			fmt.Sprintf("DROP USER IF EXISTS '%s'@'localhost';", db.Username))
		return cmd.Run()
		
	case "postgresql":
		// Drop database
		cmd := exec.Command("sudo", "-u", "postgres", "dropdb", db.Name)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to drop PostgreSQL database: %v", err)
		}
		
		// Drop user
		cmd = exec.Command("sudo", "-u", "postgres", "dropuser", db.Username)
		return cmd.Run()
		
	default:
		return fmt.Errorf("unsupported database type: %s", db.Type)
	}
}
