
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
