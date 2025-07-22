
package services

import (
	"AdminiSoftware/internal/models"
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupService struct{}

func NewBackupService() *BackupService {
	return &BackupService{}
}

func (s *BackupService) CreateBackup(userID uint, backupType string, includes []string) (*models.Backup, error) {
	backup := &models.Backup{
		UserID:     userID,
		Type:       backupType,
		Status:     "creating",
		Size:       0,
		CreatedAt:  time.Now(),
	}

	// Generate backup filename
	timestamp := time.Now().Format("20060102_150405")
	backup.Filename = fmt.Sprintf("backup_%s_%s.tar.gz", backupType, timestamp)
	backup.Path = filepath.Join("/backups", fmt.Sprintf("user_%d", userID), backup.Filename)

	return backup, nil
}

func (s *BackupService) CreateFullBackup(userID uint, userPath string) error {
	backupDir := filepath.Join("/backups", fmt.Sprintf("user_%d", userID))
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	backupFile := filepath.Join(backupDir, fmt.Sprintf("full_backup_%s.tar.gz", timestamp))

	file, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %v", err)
	}
	defer file.Close()

	gzWriter := gzip.NewWriter(file)
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return s.addToTar(tarWriter, userPath, "")
}

func (s *BackupService) addToTar(tarWriter *tar.Writer, sourcePath, basePath string) error {
	return filepath.Walk(sourcePath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		if basePath != "" {
			header.Name = filepath.Join(basePath, strings.TrimPrefix(file, sourcePath))
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !fi.IsDir() {
			data, err := os.Open(file)
			if err != nil {
				return err
			}
			defer data.Close()

			_, err = io.Copy(tarWriter, data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *BackupService) RestoreBackup(userID uint, backupID uint, restorePath string) error {
	// Implementation for restoring backup
	// This would extract the backup file to the specified path
	return fmt.Errorf("restore functionality not yet implemented")
}

func (s *BackupService) ScheduleBackup(userID uint, schedule string, backupType string) (*models.BackupSchedule, error) {
	backupSchedule := &models.BackupSchedule{
		UserID:    userID,
		Schedule:  schedule,
		Type:      backupType,
		Status:    "active",
		CreatedAt: time.Now(),
	}

	return backupSchedule, nil
}

func (s *BackupService) DeleteOldBackups(userID uint, retentionDays int) error {
	backupDir := filepath.Join("/backups", fmt.Sprintf("user_%d", userID))
	
	cutoff := time.Now().AddDate(0, 0, -retentionDays)
	
	return filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && info.ModTime().Before(cutoff) {
			return os.Remove(path)
		}
		
		return nil
	})
}

func (s *BackupService) GetBackupSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	
	return info.Size(), nil
}

func (s *BackupService) ValidateBackup(filePath string) error {
	// Open and validate the backup file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open backup file: %v", err)
	}
	defer file.Close()

	// Check if it's a valid gzip file
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("invalid gzip format: %v", err)
	}
	defer gzReader.Close()

	// Check if it's a valid tar file
	tarReader := tar.NewReader(gzReader)
	_, err = tarReader.Next()
	if err != nil && err != io.EOF {
		return fmt.Errorf("invalid tar format: %v", err)
	}

	return nil
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

type BackupService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewBackupService(db *gorm.DB, logger *utils.Logger) *BackupService {
	return &BackupService{
		db:     db,
		logger: logger,
	}
}

func (s *BackupService) CreateBackup(userID uint, req *models.CreateBackupRequest) (*models.Backup, error) {
	backup := &models.Backup{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Status:      "creating",
		Size:        0,
	}

	if err := s.db.Create(backup).Error; err != nil {
		return nil, fmt.Errorf("failed to create backup record: %v", err)
	}

	// Create backup asynchronously
	go s.performBackup(backup)

	return backup, nil
}

func (s *BackupService) performBackup(backup *models.Backup) {
	backupDir := fmt.Sprintf("/var/backups/users/%d", backup.UserID)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		s.updateBackupStatus(backup.ID, "failed", 0, err.Error())
		return
	}

	filename := fmt.Sprintf("%s_%d_%s.tar.gz", backup.Name, backup.UserID, time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(backupDir, filename)

	var cmd *exec.Cmd
	switch backup.Type {
	case "full":
		userDir := fmt.Sprintf("/home/users/%d", backup.UserID)
		cmd = exec.Command("tar", "-czf", backupPath, "-C", userDir, ".")
	case "database":
		// Backup databases only
		cmd = exec.Command("sh", "-c", fmt.Sprintf("mysqldump --all-databases > %s", backupPath))
	case "files":
		userDir := fmt.Sprintf("/home/users/%d/public_html", backup.UserID)
		cmd = exec.Command("tar", "-czf", backupPath, "-C", userDir, ".")
	default:
		s.updateBackupStatus(backup.ID, "failed", 0, "Invalid backup type")
		return
	}

	if err := cmd.Run(); err != nil {
		s.updateBackupStatus(backup.ID, "failed", 0, err.Error())
		return
	}

	// Get file size
	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		s.updateBackupStatus(backup.ID, "failed", 0, err.Error())
		return
	}

	// Update backup record
	s.updateBackupStatus(backup.ID, "completed", fileInfo.Size(), "")
	
	// Update path
	s.db.Model(&models.Backup{}).Where("id = ?", backup.ID).Update("path", backupPath)
}

func (s *BackupService) updateBackupStatus(backupID uint, status string, size int64, errorMsg string) {
	updates := map[string]interface{}{
		"status": status,
		"size":   size,
	}
	
	if errorMsg != "" {
		updates["error_message"] = errorMsg
	}

	s.db.Model(&models.Backup{}).Where("id = ?", backupID).Updates(updates)
}

func (s *BackupService) GetBackups(userID uint, page, limit int) ([]models.Backup, int64, error) {
	var backups []models.Backup
	var total int64

	offset := (page - 1) * limit

	query := s.db.Model(&models.Backup{}).Where("user_id = ?", userID)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&backups).Error; err != nil {
		return nil, 0, err
	}

	return backups, total, nil
}

func (s *BackupService) GetBackup(id uint, userID uint) (*models.Backup, error) {
	var backup models.Backup
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&backup).Error; err != nil {
		return nil, err
	}
	return &backup, nil
}

func (s *BackupService) DeleteBackup(id uint, userID uint) error {
	var backup models.Backup
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&backup).Error; err != nil {
		return err
	}

	// Delete backup file
	if backup.Path != "" {
		if err := os.Remove(backup.Path); err != nil {
			s.logger.Error(fmt.Sprintf("Failed to delete backup file %s: %v", backup.Path, err))
		}
	}

	// Delete database record
	if err := s.db.Delete(&backup).Error; err != nil {
		return fmt.Errorf("failed to delete backup record: %v", err)
	}

	return nil
}

func (s *BackupService) RestoreBackup(id uint, userID uint) error {
	var backup models.Backup
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&backup).Error; err != nil {
		return err
	}

	if backup.Status != "completed" {
		return fmt.Errorf("backup is not completed")
	}

	// Create restore process
	go s.performRestore(&backup)

	return nil
}

func (s *BackupService) performRestore(backup *models.Backup) {
	userDir := fmt.Sprintf("/home/users/%d", backup.UserID)
	
	var cmd *exec.Cmd
	switch backup.Type {
	case "full":
		cmd = exec.Command("tar", "-xzf", backup.Path, "-C", userDir)
	case "files":
		publicDir := filepath.Join(userDir, "public_html")
		cmd = exec.Command("tar", "-xzf", backup.Path, "-C", publicDir)
	}

	if err := cmd.Run(); err != nil {
		s.logger.Error(fmt.Sprintf("Failed to restore backup %d: %v", backup.ID, err))
		return
	}

	s.logger.Info(fmt.Sprintf("Backup %d restored successfully for user %d", backup.ID, backup.UserID))
}
