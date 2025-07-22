
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
