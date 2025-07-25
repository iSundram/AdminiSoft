
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
