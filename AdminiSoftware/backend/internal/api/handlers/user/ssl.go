
package user

import (
	"AdminiSoftware/internal/models"
	"AdminiSoftware/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserSSLHandler struct {
	db         *gorm.DB
	sslService *services.SSLService
}

func NewUserSSLHandler(db *gorm.DB, sslService *services.SSLService) *UserSSLHandler {
	return &UserSSLHandler{
		db:         db,
		sslService: sslService,
	}
}

func (h *UserSSLHandler) GetSSLCertificates(c *gin.Context) {
	userID := c.GetUint("user_id")

	var certificates []models.SSL
	if err := h.db.Where("user_id = ?", userID).Find(&certificates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SSL certificates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"certificates": certificates})
}

func (h *UserSSLHandler) GenerateSSLCertificate(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		Domain      string `json:"domain" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Email       string `json:"email" binding:"required"`
		Country     string `json:"country"`
		State       string `json:"state"`
		City        string `json:"city"`
		Organization string `json:"organization"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate SSL certificate
	cert, key, err := h.sslService.GenerateCertificate(req.Domain, req.Email, req.Country, req.State, req.City, req.Organization)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate SSL certificate"})
		return
	}

	ssl := models.SSL{
		UserID:      userID,
		DomainID:    req.DomainID,
		Domain:      req.Domain,
		Type:        req.Type,
		Certificate: cert,
		PrivateKey:  key,
		Status:      "active",
	}

	if err := h.db.Create(&ssl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate generated successfully", "ssl": ssl})
}

func (h *UserSSLHandler) InstallSSLCertificate(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		Domain      string `json:"domain" binding:"required"`
		Certificate string `json:"certificate" binding:"required"`
		PrivateKey  string `json:"private_key" binding:"required"`
		CABundle    string `json:"ca_bundle"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ssl := models.SSL{
		UserID:      userID,
		DomainID:    req.DomainID,
		Domain:      req.Domain,
		Type:        "manual",
		Certificate: req.Certificate,
		PrivateKey:  req.PrivateKey,
		CABundle:    req.CABundle,
		Status:      "active",
	}

	if err := h.db.Create(&ssl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate installed successfully", "ssl": ssl})
}

func (h *UserSSLHandler) RenewSSLCertificate(c *gin.Context) {
	userID := c.GetUint("user_id")
	sslID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var ssl models.SSL
	if err := h.db.Where("id = ? AND user_id = ?", sslID, userID).First(&ssl).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	// Renew SSL certificate
	cert, key, err := h.sslService.RenewCertificate(ssl.Domain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to renew SSL certificate"})
		return
	}

	ssl.Certificate = cert
	ssl.PrivateKey = key

	if err := h.db.Save(&ssl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save renewed SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate renewed successfully", "ssl": ssl})
}

func (h *UserSSLHandler) DeleteSSLCertificate(c *gin.Context) {
	userID := c.GetUint("user_id")
	sslID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var ssl models.SSL
	if err := h.db.Where("id = ? AND user_id = ?", sslID, userID).First(&ssl).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SSL certificate not found"})
		return
	}

	if err := h.db.Delete(&ssl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SSL certificate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "SSL certificate deleted successfully"})
}
