
package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileService struct {
	basePath string
}

func NewFileService(basePath string) *FileService {
	return &FileService{basePath: basePath}
}

func (s *FileService) ListFiles(userPath string) ([]FileInfo, error) {
	fullPath := filepath.Join(s.basePath, userPath)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return nil, fmt.Errorf("invalid path")
	}
	
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}
	
	var fileInfos []FileInfo
	for _, file := range files {
		info, _ := file.Info()
		fileInfos = append(fileInfos, FileInfo{
			Name:    file.Name(),
			Size:    info.Size(),
			IsDir:   file.IsDir(),
			ModTime: info.ModTime(),
		})
	}
	
	return fileInfos, nil
}

func (s *FileService) CreateDirectory(userPath, dirName string) error {
	fullPath := filepath.Join(s.basePath, userPath, dirName)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return fmt.Errorf("invalid path")
	}
	
	return os.MkdirAll(fullPath, 0755)
}

func (s *FileService) DeleteFile(userPath string) error {
	fullPath := filepath.Join(s.basePath, userPath)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return fmt.Errorf("invalid path")
	}
	
	return os.RemoveAll(fullPath)
}

func (s *FileService) ReadFile(userPath string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, userPath)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return nil, fmt.Errorf("invalid path")
	}
	
	return os.ReadFile(fullPath)
}

func (s *FileService) WriteFile(userPath string, content []byte) error {
	fullPath := filepath.Join(s.basePath, userPath)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return fmt.Errorf("invalid path")
	}
	
	return os.WriteFile(fullPath, content, 0644)
}

func (s *FileService) UploadFile(userPath string, file io.Reader) error {
	fullPath := filepath.Join(s.basePath, userPath)
	
	if !strings.HasPrefix(fullPath, s.basePath) {
		return fmt.Errorf("invalid path")
	}
	
	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, file)
	return err
}

type FileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"is_dir"`
	ModTime time.Time `json:"mod_time"`
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

type FileService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewFileService(db *gorm.DB, logger *utils.Logger) *FileService {
	return &FileService{
		db:     db,
		logger: logger,
	}
}

func (s *FileService) ListFiles(userID uint, path string) ([]models.FileInfo, error) {
	// Get user's home directory
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(user.HomeDirectory, path)
	
	// Security check - ensure path is within user's home directory
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return nil, errors.New("access denied: path outside user directory")
	}

	files, err := os.ReadDir(fullPath)
	if err != nil {
		s.logger.Error("Failed to list directory", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return nil, fmt.Errorf("failed to read directory: %v", err)
	}

	var fileInfos []models.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}

		fileInfo := models.FileInfo{
			Name:         file.Name(),
			Size:         info.Size(),
			ModifiedTime: info.ModTime(),
			IsDirectory:  file.IsDir(),
			Permissions:  info.Mode().String(),
		}

		// Get MIME type for files
		if !file.IsDir() {
			fileInfo.MimeType = s.getMimeType(file.Name())
		}

		fileInfos = append(fileInfos, fileInfo)
	}

	return fileInfos, nil
}

func (s *FileService) CreateDirectory(userID uint, path string) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(user.HomeDirectory, path)
	
	// Security check
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return errors.New("access denied: path outside user directory")
	}

	if err := os.MkdirAll(fullPath, 0755); err != nil {
		s.logger.Error("Failed to create directory", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return fmt.Errorf("failed to create directory: %v", err)
	}

	return nil
}

func (s *FileService) DeleteFile(userID uint, path string) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(user.HomeDirectory, path)
	
	// Security check
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return errors.New("access denied: path outside user directory")
	}

	if err := os.RemoveAll(fullPath); err != nil {
		s.logger.Error("Failed to delete file", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func (s *FileService) UploadFile(userID uint, path string, filename string, content io.Reader) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(user.HomeDirectory, path, filename)
	
	// Security check
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return errors.New("access denied: path outside user directory")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		s.logger.Error("Failed to create file", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, content); err != nil {
		s.logger.Error("Failed to write file content", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func (s *FileService) ReadFile(userID uint, path string) ([]byte, error) {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(user.HomeDirectory, path)
	
	// Security check
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return nil, errors.New("access denied: path outside user directory")
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		s.logger.Error("Failed to read file", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	return content, nil
}

func (s *FileService) WriteFile(userID uint, path string, content []byte) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(user.HomeDirectory, path)
	
	// Security check
	if !strings.HasPrefix(fullPath, user.HomeDirectory) {
		return errors.New("access denied: path outside user directory")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		s.logger.Error("Failed to write file", map[string]interface{}{
			"error": err.Error(),
			"path": fullPath,
		})
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}

func (s *FileService) MoveFile(userID uint, sourcePath, destPath string) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullSourcePath := filepath.Join(user.HomeDirectory, sourcePath)
	fullDestPath := filepath.Join(user.HomeDirectory, destPath)
	
	// Security checks
	if !strings.HasPrefix(fullSourcePath, user.HomeDirectory) {
		return errors.New("access denied: source path outside user directory")
	}
	if !strings.HasPrefix(fullDestPath, user.HomeDirectory) {
		return errors.New("access denied: destination path outside user directory")
	}

	if err := os.Rename(fullSourcePath, fullDestPath); err != nil {
		s.logger.Error("Failed to move file", map[string]interface{}{
			"error": err.Error(),
			"source": fullSourcePath,
			"dest": fullDestPath,
		})
		return fmt.Errorf("failed to move file: %v", err)
	}

	return nil
}

func (s *FileService) CopyFile(userID uint, sourcePath, destPath string) error {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return err
	}

	fullSourcePath := filepath.Join(user.HomeDirectory, sourcePath)
	fullDestPath := filepath.Join(user.HomeDirectory, destPath)
	
	// Security checks
	if !strings.HasPrefix(fullSourcePath, user.HomeDirectory) {
		return errors.New("access denied: source path outside user directory")
	}
	if !strings.HasPrefix(fullDestPath, user.HomeDirectory) {
		return errors.New("access denied: destination path outside user directory")
	}

	sourceFile, err := os.Open(fullSourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(fullDestPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		s.logger.Error("Failed to copy file", map[string]interface{}{
			"error": err.Error(),
			"source": fullSourcePath,
			"dest": fullDestPath,
		})
		return fmt.Errorf("failed to copy file: %v", err)
	}

	return nil
}

func (s *FileService) getUserHomeDir(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	if user.HomeDirectory == "" {
		user.HomeDirectory = fmt.Sprintf("/home/%s", user.Username)
	}

	return &user, nil
}

func (s *FileService) getMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeTypes := map[string]string{
		".txt":  "text/plain",
		".html": "text/html",
		".css":  "text/css",
		".js":   "text/javascript",
		".json": "application/json",
		".xml":  "text/xml",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".pdf":  "application/pdf",
		".zip":  "application/zip",
		".tar":  "application/tar",
		".gz":   "application/gzip",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}

	return "application/octet-stream"
}

func (s *FileService) GetDiskUsage(userID uint) (*models.DiskUsage, error) {
	user, err := s.getUserHomeDir(userID)
	if err != nil {
		return nil, err
	}

	var totalSize int64
	err = filepath.Walk(user.HomeDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files that can't be accessed
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to calculate disk usage: %v", err)
	}

	return &models.DiskUsage{
		UsedBytes:  totalSize,
		UsedMB:     totalSize / (1024 * 1024),
		UsedGB:     float64(totalSize) / (1024 * 1024 * 1024),
	}, nil
}
