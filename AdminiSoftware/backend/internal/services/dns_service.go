
package services

import (
	"AdminiSoftware/internal/models"
	"fmt"
	"net"
)

type DNSService struct{}

func NewDNSService() *DNSService {
	return &DNSService{}
}

func (s *DNSService) CreateDNSZone(domain string, userID uint) (*models.DNS, error) {
	// Validate domain
	if !s.isValidDomain(domain) {
		return nil, fmt.Errorf("invalid domain format")
	}

	dns := &models.DNS{
		UserID: userID,
		Domain: domain,
		Type:   "A",
		Name:   "@",
		Value:  "192.168.1.1", // Default IP
		TTL:    3600,
		Status: "active",
	}

	return dns, nil
}

func (s *DNSService) AddDNSRecord(domain, recordType, name, value string, ttl int, userID uint) (*models.DNS, error) {
	// Validate record type
	validTypes := []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "PTR", "SRV"}
	if !s.contains(validTypes, recordType) {
		return nil, fmt.Errorf("invalid DNS record type")
	}

	// Validate based on record type
	if err := s.validateDNSRecord(recordType, value); err != nil {
		return nil, err
	}

	dns := &models.DNS{
		UserID: userID,
		Domain: domain,
		Type:   recordType,
		Name:   name,
		Value:  value,
		TTL:    ttl,
		Status: "active",
	}

	return dns, nil
}

func (s *DNSService) UpdateDNSRecord(dns *models.DNS, value string, ttl int) error {
	// Validate the new value
	if err := s.validateDNSRecord(dns.Type, value); err != nil {
		return err
	}

	dns.Value = value
	dns.TTL = ttl

	return nil
}

func (s *DNSService) DeleteDNSRecord(dns *models.DNS) error {
	// Check if it's a critical record
	if dns.Type == "A" && dns.Name == "@" {
		return fmt.Errorf("cannot delete main A record")
	}

	return nil
}

func (s *DNSService) ValidateDNSConfiguration(domain string) error {
	// Check if domain resolves
	_, err := net.LookupHost(domain)
	if err != nil {
		return fmt.Errorf("domain does not resolve: %v", err)
	}

	return nil
}

func (s *DNSService) isValidDomain(domain string) bool {
	// Simple domain validation
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	
	// Check for valid characters and format
	// This is a simplified validation
	return true
}

func (s *DNSService) validateDNSRecord(recordType, value string) error {
	switch recordType {
	case "A":
		if net.ParseIP(value).To4() == nil {
			return fmt.Errorf("invalid IPv4 address")
		}
	case "AAAA":
		if net.ParseIP(value).To16() == nil {
			return fmt.Errorf("invalid IPv6 address")
		}
	case "CNAME":
		if !s.isValidDomain(value) {
			return fmt.Errorf("invalid domain for CNAME record")
		}
	case "MX":
		// MX records should have priority and domain
		// Format: "10 mail.example.com"
		// This is simplified validation
		if len(value) == 0 {
			return fmt.Errorf("invalid MX record format")
		}
	case "TXT":
		// TXT records can contain any text
		// No specific validation needed
	case "NS":
		if !s.isValidDomain(value) {
			return fmt.Errorf("invalid nameserver domain")
		}
	}

	return nil
}

func (s *DNSService) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
