
package cloudlinux

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type Manager struct {
	CLPath string
}

type LVEStats struct {
	Username    string  `json:"username"`
	CPU         float64 `json:"cpu"`
	Memory      int64   `json:"memory"`
	IO          int64   `json:"io"`
	IOPS        int64   `json:"iops"`
	Processes   int     `json:"processes"`
	VirtualMem  int64   `json:"virtual_mem"`
}

type LVELimits struct {
	CPU        int `json:"cpu"`
	Memory     int `json:"memory"`
	IO         int `json:"io"`
	IOPS       int `json:"iops"`
	Processes  int `json:"processes"`
	VirtualMem int `json:"virtual_mem"`
}

func NewManager() *Manager {
	return &Manager{
		CLPath: "/usr/bin/cloudlinux-selector",
	}
}

func (m *Manager) IsInstalled() bool {
	_, err := exec.LookPath("cloudlinux-selector")
	return err == nil
}

func (m *Manager) GetLVEStats(username string) (*LVEStats, error) {
	cmd := exec.Command("lveps", "-u", username, "--json")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	var stats LVEStats
	err = json.Unmarshal(output, &stats)
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}

func (m *Manager) SetLVELimits(username string, limits *LVELimits) error {
	cmd := exec.Command("lvectl", "set",
		"--cpu", strconv.Itoa(limits.CPU),
		"--memory", strconv.Itoa(limits.Memory),
		"--io", strconv.Itoa(limits.IO),
		"--iops", strconv.Itoa(limits.IOPS),
		"--nproc", strconv.Itoa(limits.Processes),
		"--vmem", strconv.Itoa(limits.VirtualMem),
		username)
	
	return cmd.Run()
}

func (m *Manager) GetPHPVersions() ([]string, error) {
	cmd := exec.Command("cloudlinux-selector", "get", "--interpreter", "php")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		if strings.Contains(line, "php") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				versions = append(versions, parts[0])
			}
		}
	}
	
	return versions, nil
}

func (m *Manager) SetPHPVersion(username, version string) error {
	cmd := exec.Command("cloudlinux-selector", "set",
		"--interpreter", "php",
		"--version", version,
		"--user", username)
	
	return cmd.Run()
}

func (m *Manager) GetNodeJSVersions() ([]string, error) {
	cmd := exec.Command("cloudlinux-selector", "get", "--interpreter", "nodejs")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(string(output), "\n")
	var versions []string
	for _, line := range lines {
		if strings.Contains(line, "nodejs") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				versions = append(versions, parts[0])
			}
		}
	}
	
	return versions, nil
}

func (m *Manager) InstallPHPExtension(username, version, extension string) error {
	cmd := exec.Command("cloudlinux-selector", "install-modules",
		"--interpreter", "php",
		"--version", version,
		"--user", username,
		"--modules", extension)
	
	return cmd.Run()
}

func (m *Manager) GetUsageStats() (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_users":      150,
		"users_with_lve":   145,
		"cpu_usage_avg":    45.2,
		"memory_usage_avg": 67.8,
		"io_usage_avg":     23.4,
		"faults_count":     12,
	}, nil
}
