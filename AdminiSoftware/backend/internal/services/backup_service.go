
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
