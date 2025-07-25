
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
