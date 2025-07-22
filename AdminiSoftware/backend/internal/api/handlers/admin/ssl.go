
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SSLHandler struct {
	db *gorm.DB
}

func NewSSLHandler(db *gorm.DB) *SSLHandler {
	return &SSLHandler{db: db}
}

func (h *SSLHandler) GetCertificates(c *gin.Context) {
	var certificates []models.SSLCertificate
	if err := h.db.Preload("User").Find(&certificates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SSL certificates"})
		return
	}
	c.JSON(http.StatusOK, certificates)
}

func (h *SSLHandler) CreateCertificate(c *gin.Context) {
	var cert models.SSLCertificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cert.Status = "pending"
	if err := h.db.Create(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SSL certificate"})
		return
	}

	c.JSON(http.StatusCreated, cert)
}

func (h *SSLHandler) InstallCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cert models.SSLCertificate
	if err := h.db.First(&cert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	cert.Status = "installed"
	if err := h.db.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *SSLHandler) RenewCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cert models.SSLCertificate
	if err := h.db.First(&cert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	cert.Status = "renewing"
	if err := h.db.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renew SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *SSLHandler) DeleteCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.SSLCertificate{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SSL certificate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate deleted successfully"})
}
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SSLHandler struct {
	db *gorm.DB
}

func NewSSLHandler(db *gorm.DB) *SSLHandler {
	return &SSLHandler{db: db}
}

func (h *SSLHandler) GetSSLCertificates(c *gin.Context) {
	var certificates []models.SSLCertificate
	if err := h.db.Preload("Domain").Find(&certificates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SSL certificates"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ssl_certificates": certificates})
}

func (h *SSLHandler) CreateSSLCertificate(c *gin.Context) {
	var cert models.SSLCertificate
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create SSL certificate"})
		return
	}

	c.JSON(http.StatusCreated, cert)
}

func (h *SSLHandler) UpdateSSLCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cert models.SSLCertificate
	if err := h.db.First(&cert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, cert)
}

func (h *SSLHandler) DeleteSSLCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.SSLCertificate{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SSL certificate"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate deleted successfully"})
}

func (h *SSLHandler) InstallSSLCertificate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var cert models.SSLCertificate
	if err := h.db.First(&cert, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	// Install SSL certificate logic here
	cert.Status = "installed"
	if err := h.db.Save(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate installed successfully"})
}

func (h *SSLHandler) GenerateLetsEncryptSSL(c *gin.Context) {
	var request struct {
		Domain string `json:"domain" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate Let's Encrypt SSL certificate
	cert := models.SSLCertificate{
		DomainID:    1, // This should be determined from the domain
		Type:        "letsencrypt",
		Status:      "active",
		Certificate: "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----",
		PrivateKey:  "-----BEGIN PRIVATE KEY-----\n...\n-----END PRIVATE KEY-----",
	}

	if err := h.db.Create(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate SSL certificate"})
		return
	}

	c.JSON(http.StatusCreated, cert)
}
