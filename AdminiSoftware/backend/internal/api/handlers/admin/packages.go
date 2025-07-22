
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PackageHandler struct {
	db *gorm.DB
}

func NewPackageHandler(db *gorm.DB) *PackageHandler {
	return &PackageHandler{db: db}
}

func (h *PackageHandler) GetPackages(c *gin.Context) {
	var packages []models.Package
	if err := h.db.Find(&packages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch packages"})
		return
	}
	c.JSON(http.StatusOK, packages)
}

func (h *PackageHandler) CreatePackage(c *gin.Context) {
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create package"})
		return
	}

	c.JSON(http.StatusCreated, pkg)
}

func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var pkg models.Package
	if err := h.db.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Package not found"})
		return
	}

	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update package"})
		return
	}

	c.JSON(http.StatusOK, pkg)
}

func (h *PackageHandler) DeletePackage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.Package{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete package"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Package deleted successfully"})
}

func (h *PackageHandler) GetFeatures(c *gin.Context) {
	features := map[string]interface{}{
		"disk_space":          []string{"100MB", "500MB", "1GB", "5GB", "10GB", "50GB", "100GB", "unlimited"},
		"bandwidth":           []string{"1GB", "5GB", "10GB", "50GB", "100GB", "500GB", "unlimited"},
		"email_accounts":      []string{"1", "5", "10", "25", "50", "100", "unlimited"},
		"databases":           []string{"1", "5", "10", "25", "50", "unlimited"},
		"subdomains":          []string{"1", "5", "10", "25", "50", "unlimited"},
		"addon_domains":       []string{"0", "1", "5", "10", "25", "unlimited"},
		"parked_domains":      []string{"0", "1", "5", "10", "25", "unlimited"},
		"ftp_accounts":        []string{"1", "5", "10", "25", "50", "unlimited"},
		"cron_jobs":           []string{"0", "1", "5", "10", "25", "unlimited"},
		"backup_enabled":      []bool{true, false},
		"ssl_enabled":         []bool{true, false},
		"dedicated_ip":        []bool{true, false},
		"shell_access":        []bool{true, false},
		"max_email_per_hour":  []string{"50", "100", "250", "500", "1000", "unlimited"},
		"max_email_accounts":  []string{"1", "10", "25", "50", "100", "unlimited"},
		"max_mailing_lists":   []string{"0", "1", "5", "10", "25", "unlimited"},
		"max_autoresponders":  []string{"0", "1", "5", "10", "25", "unlimited"},
		"max_forwarders":      []string{"0", "1", "5", "10", "25", "unlimited"},
	}
	c.JSON(http.StatusOK, features)
}
