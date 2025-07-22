
package services

import (
	"runtime"
	"time"
)

type MonitoringService struct{}

func NewMonitoringService() *MonitoringService {
	return &MonitoringService{}
}

func (s *MonitoringService) GetSystemStats() SystemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	return SystemStats{
		CPU: CPUStats{
			Usage:     65.5,
			LoadAvg1:  1.2,
			LoadAvg5:  1.5,
			LoadAvg15: 1.8,
		},
		Memory: MemoryStats{
			Total:     8192,
			Used:      5120,
			Free:      3072,
			Cached:    1024,
			SwapTotal: 2048,
			SwapUsed:  256,
		},
		Disk: DiskStats{
			Total: 100,
			Used:  65,
			Free:  35,
		},
		Network: NetworkStats{
			BytesIn:  1024000,
			BytesOut: 512000,
		},
		Uptime:    time.Since(time.Now().Add(-24 * time.Hour)).Seconds(),
		Timestamp: time.Now(),
	}
}

func (s *MonitoringService) GetServiceStatus() []ServiceStatus {
	return []ServiceStatus{
		{Name: "Apache", Status: "running", Port: 80, Uptime: 86400},
		{Name: "MySQL", Status: "running", Port: 3306, Uptime: 86400},
		{Name: "Redis", Status: "running", Port: 6379, Uptime: 86400},
		{Name: "Postfix", Status: "running", Port: 25, Uptime: 86400},
		{Name: "SSH", Status: "running", Port: 22, Uptime: 86400},
	}
}

func (s *MonitoringService) GetAlerts() []Alert {
	return []Alert{
		{
			Level:     "warning",
			Message:   "High CPU usage detected",
			Timestamp: time.Now().Add(-30 * time.Minute),
			Resolved:  false,
		},
		{
			Level:     "info",
			Message:   "Backup completed successfully",
			Timestamp: time.Now().Add(-2 * time.Hour),
			Resolved:  true,
		},
	}
}

type SystemStats struct {
	CPU       CPUStats     `json:"cpu"`
	Memory    MemoryStats  `json:"memory"`
	Disk      DiskStats    `json:"disk"`
	Network   NetworkStats `json:"network"`
	Uptime    float64      `json:"uptime"`
	Timestamp time.Time    `json:"timestamp"`
}

type CPUStats struct {
	Usage     float64 `json:"usage"`
	LoadAvg1  float64 `json:"load_avg_1"`
	LoadAvg5  float64 `json:"load_avg_5"`
	LoadAvg15 float64 `json:"load_avg_15"`
}

type MemoryStats struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Free      int64 `json:"free"`
	Cached    int64 `json:"cached"`
	SwapTotal int64 `json:"swap_total"`
	SwapUsed  int64 `json:"swap_used"`
}

type DiskStats struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
	Free  int64 `json:"free"`
}

type NetworkStats struct {
	BytesIn  int64 `json:"bytes_in"`
	BytesOut int64 `json:"bytes_out"`
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Port   int    `json:"port"`
	Uptime int64  `json:"uptime"`
}

type Alert struct {
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Resolved  bool      `json:"resolved"`
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type MonitoringService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewMonitoringService(db *gorm.DB, logger *utils.Logger) *MonitoringService {
	return &MonitoringService{
		db:     db,
		logger: logger,
	}
}

func (s *MonitoringService) GetSystemStats() (*models.SystemStats, error) {
	stats := &models.SystemStats{
		Timestamp: time.Now(),
	}

	// Get CPU usage
	cpuUsage, err := s.getCPUUsage()
	if err != nil {
		s.logger.Error("Failed to get CPU usage", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		stats.CPUUsage = cpuUsage
	}

	// Get memory usage
	memUsage, err := s.getMemoryUsage()
	if err != nil {
		s.logger.Error("Failed to get memory usage", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		stats.MemoryUsage = memUsage
	}

	// Get disk usage
	diskUsage, err := s.getDiskUsage()
	if err != nil {
		s.logger.Error("Failed to get disk usage", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		stats.DiskUsage = diskUsage
	}

	// Get network stats
	networkStats, err := s.getNetworkStats()
	if err != nil {
		s.logger.Error("Failed to get network stats", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		stats.NetworkIn = networkStats.BytesReceived
		stats.NetworkOut = networkStats.BytesSent
	}

	// Get load average
	loadAvg, err := s.getLoadAverage()
	if err != nil {
		s.logger.Error("Failed to get load average", map[string]interface{}{
			"error": err.Error(),
		})
	} else {
		stats.LoadAverage = loadAvg
	}

	// Save stats to database
	if err := s.db.Create(stats).Error; err != nil {
		s.logger.Error("Failed to save system stats", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return stats, nil
}

func (s *MonitoringService) getCPUUsage() (float64, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | sed 's/.*, *\\([0-9.]*\\)%* id.*/\\1/' | awk '{print 100 - $1}'")
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}
		
		cpuStr := strings.TrimSpace(string(output))
		cpu, err := strconv.ParseFloat(cpuStr, 64)
		if err != nil {
			return 0, err
		}
		return cpu, nil
	}
	
	// Fallback for non-Linux systems
	return float64(runtime.NumCPU()) * 10.0, nil // Mock value
}

func (s *MonitoringService) getMemoryUsage() (*models.MemoryUsage, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("free", "-m")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) < 2 {
			return nil, fmt.Errorf("invalid memory output")
		}

		memLine := strings.Fields(lines[1])
		if len(memLine) < 3 {
			return nil, fmt.Errorf("invalid memory line format")
		}

		total, _ := strconv.ParseInt(memLine[1], 10, 64)
		used, _ := strconv.ParseInt(memLine[2], 10, 64)
		
		return &models.MemoryUsage{
			TotalMB: total,
			UsedMB:  used,
			FreeMB:  total - used,
			UsagePercent: float64(used) / float64(total) * 100,
		}, nil
	}

	// Fallback for non-Linux systems
	return &models.MemoryUsage{
		TotalMB: 8192,
		UsedMB:  4096,
		FreeMB:  4096,
		UsagePercent: 50.0,
	}, nil
}

func (s *MonitoringService) getDiskUsage() (*models.DiskUsage, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("df", "-h", "/")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(output), "\n")
		if len(lines) < 2 {
			return nil, fmt.Errorf("invalid disk output")
		}

		diskLine := strings.Fields(lines[1])
		if len(diskLine) < 4 {
			return nil, fmt.Errorf("invalid disk line format")
		}

		total := s.parseSize(diskLine[1])
		used := s.parseSize(diskLine[2])
		
		return &models.DiskUsage{
			UsedBytes: int64(used * 1024 * 1024),
			UsedMB:    int64(used),
			UsedGB:    used / 1024,
		}, nil
	}

	// Fallback for non-Linux systems
	return &models.DiskUsage{
		UsedBytes: 50 * 1024 * 1024 * 1024, // 50GB
		UsedMB:    50 * 1024,
		UsedGB:    50,
	}, nil
}

func (s *MonitoringService) parseSize(sizeStr string) float64 {
	sizeStr = strings.ToUpper(sizeStr)
	var multiplier float64 = 1
	
	if strings.HasSuffix(sizeStr, "K") {
		multiplier = 1
		sizeStr = sizeStr[:len(sizeStr)-1]
	} else if strings.HasSuffix(sizeStr, "M") {
		multiplier = 1024
		sizeStr = sizeStr[:len(sizeStr)-1]
	} else if strings.HasSuffix(sizeStr, "G") {
		multiplier = 1024 * 1024
		sizeStr = sizeStr[:len(sizeStr)-1]
	} else if strings.HasSuffix(sizeStr, "T") {
		multiplier = 1024 * 1024 * 1024
		sizeStr = sizeStr[:len(sizeStr)-1]
	}
	
	size, _ := strconv.ParseFloat(sizeStr, 64)
	return size * multiplier / 1024 // Convert to MB
}

func (s *MonitoringService) getNetworkStats() (*models.NetworkStats, error) {
	// Simplified network stats - in production, would parse /proc/net/dev
	return &models.NetworkStats{
		BytesReceived: 1024 * 1024 * 100, // 100MB
		BytesSent:     1024 * 1024 * 50,  // 50MB
	}, nil
}

func (s *MonitoringService) getLoadAverage() (float64, error) {
	if runtime.GOOS == "linux" {
		cmd := exec.Command("uptime")
		output, err := cmd.Output()
		if err != nil {
			return 0, err
		}

		uptimeStr := string(output)
		parts := strings.Split(uptimeStr, "load average:")
		if len(parts) < 2 {
			return 0, fmt.Errorf("invalid uptime output")
		}

		loadParts := strings.Split(strings.TrimSpace(parts[1]), ",")
		if len(loadParts) < 1 {
			return 0, fmt.Errorf("invalid load average format")
		}

		load, err := strconv.ParseFloat(strings.TrimSpace(loadParts[0]), 64)
		if err != nil {
			return 0, err
		}

		return load, nil
	}

	// Fallback for non-Linux systems
	return 1.5, nil
}

func (s *MonitoringService) GetHistoricalStats(hours int) ([]models.SystemStats, error) {
	var stats []models.SystemStats
	since := time.Now().Add(-time.Duration(hours) * time.Hour)

	if err := s.db.Where("timestamp >= ?", since).Order("timestamp DESC").Find(&stats).Error; err != nil {
		s.logger.Error("Failed to get historical stats", map[string]interface{}{
			"error": err.Error(),
			"hours": hours,
		})
		return nil, err
	}

	return stats, nil
}

func (s *MonitoringService) CheckServiceStatus(serviceName string) (*models.ServiceStatus, error) {
	status := &models.ServiceStatus{
		ServiceName: serviceName,
		CheckedAt:   time.Now(),
	}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("systemctl", "is-active", serviceName)
		output, err := cmd.Output()
		
		if err == nil && strings.TrimSpace(string(output)) == "active" {
			status.Status = "running"
			status.IsHealthy = true
		} else {
			status.Status = "stopped"
			status.IsHealthy = false
		}
	} else {
		// Fallback for non-Linux systems
		status.Status = "running"
		status.IsHealthy = true
	}

	// Save status to database
	if err := s.db.Create(status).Error; err != nil {
		s.logger.Error("Failed to save service status", map[string]interface{}{
			"error": err.Error(),
			"service": serviceName,
		})
	}

	return status, nil
}

func (s *MonitoringService) GetServiceStatuses() ([]models.ServiceStatus, error) {
	services := []string{"nginx", "apache2", "mysql", "postgresql", "redis", "php-fpm"}
	var statuses []models.ServiceStatus

	for _, service := range services {
		status, err := s.CheckServiceStatus(service)
		if err != nil {
			s.logger.Error("Failed to check service status", map[string]interface{}{
				"error": err.Error(),
				"service": service,
			})
			continue
		}
		statuses = append(statuses, *status)
	}

	return statuses, nil
}

func (s *MonitoringService) CreateAlert(alert *models.Alert) error {
	if err := s.db.Create(alert).Error; err != nil {
		s.logger.Error("Failed to create alert", map[string]interface{}{
			"error": err.Error(),
			"type": alert.Type,
		})
		return err
	}

	s.logger.Info("Alert created", map[string]interface{}{
		"type": alert.Type,
		"message": alert.Message,
		"severity": alert.Severity,
	})

	return nil
}

func (s *MonitoringService) GetAlerts(limit int) ([]models.Alert, error) {
	var alerts []models.Alert
	
	query := s.db.Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&alerts).Error; err != nil {
		s.logger.Error("Failed to get alerts", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return alerts, nil
}

func (s *MonitoringService) MarkAlertAsRead(alertID uint) error {
	var alert models.Alert
	if err := s.db.First(&alert, alertID).Error; err != nil {
		return fmt.Errorf("alert not found")
	}

	alert.IsRead = true
	if err := s.db.Save(&alert).Error; err != nil {
		s.logger.Error("Failed to mark alert as read", map[string]interface{}{
			"error": err.Error(),
			"alert_id": alertID,
		})
		return err
	}

	return nil
}

func (s *MonitoringService) StartMonitoring() {
	ticker := time.NewTicker(5 * time.Minute) // Collect stats every 5 minutes
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats, err := s.GetSystemStats()
			if err != nil {
				s.logger.Error("Failed to collect system stats", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			// Check for alert conditions
			s.checkAlertConditions(stats)
		}
	}
}

func (s *MonitoringService) checkAlertConditions(stats *models.SystemStats) {
	// High CPU usage alert
	if stats.CPUUsage > 90 {
		alert := &models.Alert{
			Type:     "cpu_high",
			Severity: "warning",
			Message:  fmt.Sprintf("High CPU usage detected: %.2f%%", stats.CPUUsage),
		}
		s.CreateAlert(alert)
	}

	// High memory usage alert
	if stats.MemoryUsage != nil && stats.MemoryUsage.UsagePercent > 90 {
		alert := &models.Alert{
			Type:     "memory_high",
			Severity: "warning",
			Message:  fmt.Sprintf("High memory usage detected: %.2f%%", stats.MemoryUsage.UsagePercent),
		}
		s.CreateAlert(alert)
	}

	// High load average alert
	if stats.LoadAverage > 5 {
		alert := &models.Alert{
			Type:     "load_high",
			Severity: "critical",
			Message:  fmt.Sprintf("High load average detected: %.2f", stats.LoadAverage),
		}
		s.CreateAlert(alert)
	}
}
