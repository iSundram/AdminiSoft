
package utils

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

type Validator struct{}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

func (v *Validator) ValidateDomain(domain string) error {
	if len(domain) == 0 || len(domain) > 253 {
		return fmt.Errorf("domain length must be between 1 and 253 characters")
	}

	// Remove trailing dot if present
	domain = strings.TrimSuffix(domain, ".")

	// Split domain into labels
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return fmt.Errorf("domain must have at least two labels")
	}

	// Validate each label
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 {
			return fmt.Errorf("domain label length must be between 1 and 63 characters")
		}

		// Check for valid characters
		labelRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
		if !labelRegex.MatchString(label) {
			return fmt.Errorf("invalid characters in domain label: %s", label)
		}
	}

	return nil
}

func (v *Validator) ValidateIP(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("invalid IP address format")
	}
	return nil
}

func (v *Validator) ValidatePort(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func (v *Validator) ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 32 {
		return fmt.Errorf("username must be between 3 and 32 characters")
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if !usernameRegex.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

func (v *Validator) ValidatePath(path string) error {
	if len(path) == 0 {
		return fmt.Errorf("path cannot be empty")
	}

	// Check for invalid characters
	invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range invalidChars {
		if strings.Contains(path, char) {
			return fmt.Errorf("path contains invalid character: %s", char)
		}
	}

	return nil
}

func (v *Validator) ValidateSubdomain(subdomain string) error {
	if len(subdomain) == 0 || len(subdomain) > 63 {
		return fmt.Errorf("subdomain length must be between 1 and 63 characters")
	}

	subdomainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?$`)
	if !subdomainRegex.MatchString(subdomain) {
		return fmt.Errorf("invalid subdomain format")
	}

	return nil
}

func (v *Validator) ValidateDBName(dbname string) error {
	if len(dbname) < 1 || len(dbname) > 64 {
		return fmt.Errorf("database name must be between 1 and 64 characters")
	}

	dbnameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !dbnameRegex.MatchString(dbname) {
		return fmt.Errorf("database name can only contain letters, numbers, and underscores")
	}

	return nil
}

func (v *Validator) ValidateCronSchedule(schedule string) error {
	// Basic cron validation (5 fields: minute hour day month weekday)
	fields := strings.Fields(schedule)
	if len(fields) != 5 {
		return fmt.Errorf("cron schedule must have 5 fields")
	}

	// More detailed validation could be added here
	return nil
}
