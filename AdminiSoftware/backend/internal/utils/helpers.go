
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}

// GenerateAPIKey generates a random API key
func GenerateAPIKey() string {
	return GenerateRandomString(32)
}

// FormatBytes formats bytes to human readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ParseBytes parses human readable bytes format to int64
func ParseBytes(s string) (int64, error) {
	s = strings.TrimSpace(strings.ToUpper(s))
	
	multipliers := map[string]int64{
		"B":  1,
		"KB": 1024,
		"MB": 1024 * 1024,
		"GB": 1024 * 1024 * 1024,
		"TB": 1024 * 1024 * 1024 * 1024,
	}
	
	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(s, suffix) {
			numStr := strings.TrimSuffix(s, suffix)
			num, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return 0, err
			}
			return int64(num * float64(multiplier)), nil
		}
	}
	
	// If no suffix, assume bytes
	return strconv.ParseInt(s, 10, 64)
}

// CalculatePercentage calculates percentage
func CalculatePercentage(used, total int64) float64 {
	if total == 0 {
		return 0
	}
	return math.Round((float64(used)/float64(total))*100*100) / 100
}

// EnsureDirectory ensures directory exists
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// GetFileSize returns file size in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetDirectorySize returns directory size in bytes
func GetDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size, err
}

// SanitizeFilename sanitizes filename for safe filesystem use
func SanitizeFilename(filename string) string {
	// Replace invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	sanitized := filename
	
	for _, char := range invalid {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}
	
	// Limit length
	if len(sanitized) > 255 {
		sanitized = sanitized[:255]
	}
	
	return sanitized
}

// GenerateSlug generates URL-friendly slug from string
func GenerateSlug(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	
	// Remove special characters
	var result strings.Builder
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	
	return result.String()
}

// FormatDuration formats duration to human readable format
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	} else if d < 24*time.Hour {
		return fmt.Sprintf("%.0fh", d.Hours())
	} else {
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd", days)
	}
}

// TruncateString truncates string to specified length
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// StringSliceContains checks if string slice contains a value
func StringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IntSliceContains checks if int slice contains a value
func IntSliceContains(slice []int, item int) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

// RemoveFromStringSlice removes item from string slice
func RemoveFromStringSlice(slice []string, item string) []string {
	var result []string
	for _, s := range slice {
		if s != item {
			result = append(result, s)
		}
	}
	return result
}

// GetEnvOrDefault gets environment variable or returns default value
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// IsValidJSON checks if string is valid JSON
func IsValidJSON(s string) bool {
	// Simple JSON validation - could be improved
	s = strings.TrimSpace(s)
	return (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}")) ||
		   (strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]"))
}

// MaskEmail masks email address for display
func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}
	
	local := parts[0]
	domain := parts[1]
	
	if len(local) <= 2 {
		return email
	}
	
	masked := local[:1] + strings.Repeat("*", len(local)-2) + local[len(local)-1:]
	return masked + "@" + domain
}

// TimeAgo returns human readable time ago format
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	if diff < time.Minute {
		return "just now"
	} else if diff < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	} else if diff < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	} else if diff < 7*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	} else {
		return t.Format("Jan 2, 2006")
	}
}
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.0fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1fh", d.Hours())
	}
	return fmt.Sprintf("%.1fd", d.Hours()/24)
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func SlugifyString(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return s
}

func GetClientIP(remoteAddr, xForwardedFor, xRealIP string) string {
	if xRealIP != "" {
		return xRealIP
	}
	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}
	return strings.Split(remoteAddr, ":")[0]
}
