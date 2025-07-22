
package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type SSLService struct{}

func NewSSLService() *SSLService {
	return &SSLService{}
}

func (s *SSLService) GenerateCertificate(domain, email, country, state, city, organization string) (string, string, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Country:      []string{country},
			Province:     []string{state},
			Locality:     []string{city},
			Organization: []string{organization},
			CommonName:   domain,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  nil,
		DNSNames:     []string{domain},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key to PEM
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return string(certPEM), string(privateKeyPEM), nil
}

func (s *SSLService) GenerateCSR(domain, email, country, state, city, organization string) (string, string, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create CSR template
	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:      []string{country},
			Province:     []string{state},
			Locality:     []string{city},
			Organization: []string{organization},
			CommonName:   domain,
		},
		DNSNames: []string{domain},
	}

	// Create CSR
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to create CSR: %v", err)
	}

	// Encode CSR to PEM
	csrPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrDER,
	})

	// Encode private key to PEM
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return string(csrPEM), string(privateKeyPEM), nil
}

func (s *SSLService) ValidateCertificate(certPEM string) error {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Check if certificate is expired
	if time.Now().After(cert.NotAfter) {
		return fmt.Errorf("certificate is expired")
	}

	// Check if certificate is not yet valid
	if time.Now().Before(cert.NotBefore) {
		return fmt.Errorf("certificate is not yet valid")
	}

	return nil
}

func (s *SSLService) GetCertificateInfo(certPEM string) (map[string]interface{}, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode certificate PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	info := map[string]interface{}{
		"subject":     cert.Subject.String(),
		"issuer":      cert.Issuer.String(),
		"not_before":  cert.NotBefore,
		"not_after":   cert.NotAfter,
		"dns_names":   cert.DNSNames,
		"serial":      cert.SerialNumber.String(),
		"algorithm":   cert.SignatureAlgorithm.String(),
		"is_ca":       cert.IsCA,
		"key_usage":   cert.KeyUsage,
	}

	return info, nil
}

func (s *SSLService) RenewCertificate(domain string) (string, string, error) {
	// In a real implementation, this would interact with Let's Encrypt or other CA
	// For now, we'll generate a new self-signed certificate
	return s.GenerateCertificate(domain, "", "US", "CA", "San Francisco", "AdminiSoftware")
}

func (s *SSLService) InstallLetsEncryptCertificate(domain string) error {
	// This would interact with Let's Encrypt ACME client
	// Implementation would use libraries like golang.org/x/crypto/acme
	return fmt.Errorf("Let's Encrypt integration not yet implemented")
}
package services

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/utils"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"

	"gorm.io/gorm"
)

type SSLService struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewSSLService(db *gorm.DB, logger *utils.Logger) *SSLService {
	return &SSLService{
		db:     db,
		logger: logger,
	}
}

func (s *SSLService) RequestCertificate(cert *models.SSLCertificate) error {
	if cert.Domain == "" {
		return errors.New("domain is required")
	}

	// Check if certificate already exists for domain
	var existing models.SSLCertificate
	if err := s.db.Where("domain = ? AND status = 'active'", cert.Domain).First(&existing).Error; err == nil {
		return errors.New("active certificate already exists for this domain")
	}

	// Set default values
	cert.Status = "pending"
	cert.RequestedAt = time.Now()
	cert.Type = "letsencrypt" // Default to Let's Encrypt

	if err := s.db.Create(cert).Error; err != nil {
		s.logger.Error("Failed to create SSL certificate request", map[string]interface{}{
			"error": err.Error(),
			"domain": cert.Domain,
		})
		return err
	}

	// Start certificate generation process (async)
	go s.processCertificateRequest(cert.ID)

	return nil
}

func (s *SSLService) processCertificateRequest(certID uint) {
	var cert models.SSLCertificate
	if err := s.db.First(&cert, certID).Error; err != nil {
		s.logger.Error("Certificate not found for processing", map[string]interface{}{
			"cert_id": certID,
		})
		return
	}

	// For demo purposes, generate a self-signed certificate
	// In production, this would integrate with Let's Encrypt or other CA
	if err := s.generateSelfSignedCertificate(&cert); err != nil {
		cert.Status = "failed"
		cert.ErrorMessage = err.Error()
		s.db.Save(&cert)
		return
	}

	cert.Status = "active"
	cert.IssuedAt = time.Now()
	cert.ExpiresAt = time.Now().AddDate(0, 3, 0) // 3 months validity
	s.db.Save(&cert)

	s.logger.Info("SSL certificate generated successfully", map[string]interface{}{
		"domain": cert.Domain,
		"cert_id": cert.ID,
	})
}

func (s *SSLService) generateSelfSignedCertificate(cert *models.SSLCertificate) error {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"AdminiSoftware"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
			CommonName:    cert.Domain,
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 3, 0), // 3 months
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses:  nil,
		DNSNames:     []string{cert.Domain, "www." + cert.Domain},
	}

	// Generate certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %v", err)
	}

	// Encode certificate to PEM
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key to PEM
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Update certificate record
	cert.Certificate = string(certPEM)
	cert.PrivateKey = string(privateKeyPEM)

	return nil
}

func (s *SSLService) InstallCertificate(certID uint, certificate, privateKey, chainCert string) error {
	var cert models.SSLCertificate
	if err := s.db.First(&cert, certID).Error; err != nil {
		return errors.New("certificate not found")
	}

	// Validate certificate format
	if certificate == "" || privateKey == "" {
		return errors.New("certificate and private key are required")
	}

	// Parse and validate certificate
	block, _ := pem.Decode([]byte(certificate))
	if block == nil {
		return errors.New("invalid certificate format")
	}

	parsedCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Update certificate
	cert.Certificate = certificate
	cert.PrivateKey = privateKey
	cert.ChainCertificate = chainCert
	cert.Status = "active"
	cert.IssuedAt = time.Now()
	cert.ExpiresAt = parsedCert.NotAfter
	cert.Issuer = parsedCert.Issuer.CommonName

	if err := s.db.Save(&cert).Error; err != nil {
		s.logger.Error("Failed to install SSL certificate", map[string]interface{}{
			"error": err.Error(),
			"cert_id": certID,
		})
		return err
	}

	return nil
}

func (s *SSLService) RevokeCertificate(certID uint) error {
	var cert models.SSLCertificate
	if err := s.db.First(&cert, certID).Error; err != nil {
		return errors.New("certificate not found")
	}

	cert.Status = "revoked"
	cert.RevokedAt = time.Now()

	if err := s.db.Save(&cert).Error; err != nil {
		s.logger.Error("Failed to revoke SSL certificate", map[string]interface{}{
			"error": err.Error(),
			"cert_id": certID,
		})
		return err
	}

	return nil
}

func (s *SSLService) ListCertificates(userID uint, role string) ([]models.SSLCertificate, error) {
	var certificates []models.SSLCertificate
	query := s.db

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&certificates).Error; err != nil {
		s.logger.Error("Failed to list SSL certificates", map[string]interface{}{
			"error": err.Error(),
			"user_id": userID,
		})
		return nil, err
	}

	// Remove sensitive data from response
	for i := range certificates {
		certificates[i].PrivateKey = "" // Don't include private key in list
	}

	return certificates, nil
}

func (s *SSLService) GetCertificate(certID uint) (*models.SSLCertificate, error) {
	var cert models.SSLCertificate
	if err := s.db.First(&cert, certID).Error; err != nil {
		return nil, errors.New("certificate not found")
	}

	return &cert, nil
}

func (s *SSLService) CheckExpiringCertificates(days int) ([]models.SSLCertificate, error) {
	var certificates []models.SSLCertificate
	expirationDate := time.Now().AddDate(0, 0, days)

	if err := s.db.Where("status = 'active' AND expires_at <= ?", expirationDate).Find(&certificates).Error; err != nil {
		s.logger.Error("Failed to check expiring certificates", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return certificates, nil
}

func (s *SSLService) RenewCertificate(certID uint) error {
	var cert models.SSLCertificate
	if err := s.db.First(&cert, certID).Error; err != nil {
		return errors.New("certificate not found")
	}

	if cert.Type != "letsencrypt" {
		return errors.New("only Let's Encrypt certificates can be automatically renewed")
	}

	// Create new certificate request
	newCert := models.SSLCertificate{
		Domain:      cert.Domain,
		UserID:      cert.UserID,
		Type:        cert.Type,
		Status:      "pending",
		RequestedAt: time.Now(),
	}

	if err := s.db.Create(&newCert).Error; err != nil {
		return err
	}

	// Mark old certificate as replaced
	cert.Status = "replaced"
	cert.ReplacedBy = &newCert.ID
	s.db.Save(&cert)

	// Start renewal process
	go s.processCertificateRequest(newCert.ID)

	return nil
}
