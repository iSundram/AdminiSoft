
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
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"errors"

	"gorm.io/gorm"
)

type DNSService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewDNSService(db *gorm.DB, logger *utils.Logger) *DNSService {
	return &DNSService{
		db:     db,
		logger: logger,
	}
}

func (s *DNSService) CreateZone(zone *models.DNSZone) error {
	// Validate zone
	if zone.Domain == "" {
		return errors.New("domain is required")
	}

	// Check if zone already exists
	var existingZone models.DNSZone
	if err := s.db.Where("domain = ?", zone.Domain).First(&existingZone).Error; err == nil {
		return errors.New("DNS zone already exists for this domain")
	}

	// Create zone
	if err := s.db.Create(zone).Error; err != nil {
		s.logger.Error("Failed to create DNS zone", map[string]interface{}{
			"error": err.Error(),
			"domain": zone.Domain,
		})
		return err
	}

	// Create default DNS records
	defaultRecords := []models.DNSRecord{
		{
			ZoneID: zone.ID,
			Name:   "@",
			Type:   "A",
			Value:  "127.0.0.1", // Should be replaced with actual server IP
			TTL:    3600,
		},
		{
			ZoneID: zone.ID,
			Name:   "www",
			Type:   "CNAME",
			Value:  zone.Domain,
			TTL:    3600,
		},
		{
			ZoneID: zone.ID,
			Name:   "@",
			Type:   "MX",
			Value:  "mail." + zone.Domain,
			TTL:    3600,
			Priority: 10,
		},
	}

	for _, record := range defaultRecords {
		if err := s.db.Create(&record).Error; err != nil {
			s.logger.Error("Failed to create default DNS record", map[string]interface{}{
				"error": err.Error(),
				"zone": zone.Domain,
			})
		}
	}

	return nil
}

func (s *DNSService) UpdateZone(zoneID uint, updates map[string]interface{}) error {
	var zone models.DNSZone
	if err := s.db.First(&zone, zoneID).Error; err != nil {
		return errors.New("DNS zone not found")
	}

	if err := s.db.Model(&zone).Updates(updates).Error; err != nil {
		s.logger.Error("Failed to update DNS zone", map[string]interface{}{
			"error": err.Error(),
			"zone_id": zoneID,
		})
		return err
	}

	return nil
}

func (s *DNSService) DeleteZone(zoneID uint) error {
	var zone models.DNSZone
	if err := s.db.First(&zone, zoneID).Error; err != nil {
		return errors.New("DNS zone not found")
	}

	// Delete all records first
	if err := s.db.Where("zone_id = ?", zoneID).Delete(&models.DNSRecord{}).Error; err != nil {
		s.logger.Error("Failed to delete DNS records", map[string]interface{}{
			"error": err.Error(),
			"zone_id": zoneID,
		})
		return err
	}

	// Delete zone
	if err := s.db.Delete(&zone).Error; err != nil {
		s.logger.Error("Failed to delete DNS zone", map[string]interface{}{
			"error": err.Error(),
			"zone_id": zoneID,
		})
		return err
	}

	return nil
}

func (s *DNSService) ListZones(userID uint, role string) ([]models.DNSZone, error) {
	var zones []models.DNSZone
	query := s.db.Preload("Records")

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&zones).Error; err != nil {
		s.logger.Error("Failed to list DNS zones", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		})
		return nil, err
	}

	return zones, nil
}

func (s *DNSService) CreateRecord(record *models.DNSRecord) error {
	// Validate record
	if record.ZoneID == 0 || record.Name == "" || record.Type == "" || record.Value == "" {
		return errors.New("invalid DNS record data")
	}

	// Verify zone exists
	var zone models.DNSZone
	if err := s.db.First(&zone, record.ZoneID).Error; err != nil {
		return errors.New("DNS zone not found")
	}

	if err := s.db.Create(record).Error; err != nil {
		s.logger.Error("Failed to create DNS record", map[string]interface{}{
			"error": err.Error(),
			"zone_id": record.ZoneID,
		})
		return err
	}

	return nil
}

func (s *DNSService) UpdateRecord(recordID uint, updates map[string]interface{}) error {
	var record models.DNSRecord
	if err := s.db.First(&record, recordID).Error; err != nil {
		return errors.New("DNS record not found")
	}

	if err := s.db.Model(&record).Updates(updates).Error; err != nil {
		s.logger.Error("Failed to update DNS record", map[string]interface{}{
			"error": err.Error(),
			"record_id": recordID,
		})
		return err
	}

	return nil
}

func (s *DNSService) DeleteRecord(recordID uint) error {
	var record models.DNSRecord
	if err := s.db.First(&record, recordID).Error; err != nil {
		return errors.New("DNS record not found")
	}

	if err := s.db.Delete(&record).Error; err != nil {
		s.logger.Error("Failed to delete DNS record", map[string]interface{}{
			"error": err.Error(),
			"record_id": recordID,
		})
		return err
	}

	return nil
}
