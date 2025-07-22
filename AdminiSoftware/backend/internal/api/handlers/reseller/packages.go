
package reseller

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
	resellerID := c.GetUint("user_id")
	var packages []models.Package
	if err := h.db.Where("reseller_id = ?", resellerID).Find(&packages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch packages"})
		return
	}
	c.JSON(http.StatusOK, packages)
}

func (h *PackageHandler) CreatePackage(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	var pkg models.Package
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pkg.ResellerID = &resellerID
	if err := h.db.Create(&pkg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create package"})
		return
	}

	c.JSON(http.StatusCreated, pkg)
}

func (h *PackageHandler) UpdatePackage(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var pkg models.Package
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).First(&pkg).Error; err != nil {
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
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Where("id = ? AND reseller_id = ?", id, resellerID).Delete(&models.Package{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Package deleted successfully"})
}

func (h *PackageHandler) GetPackageUsage(c *gin.Context) {
	resellerID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var accountCount int64
	h.db.Model(&models.User{}).Where("package_id = ? AND reseller_id = ?", id, resellerID).Count(&accountCount)

	usage := map[string]interface{}{
		"package_id":      id,
		"accounts_count":  accountCount,
		"disk_used":       "45.6 GB",
		"bandwidth_used":  "234.5 GB",
		"email_accounts":  156,
		"databases":       89,
		"domains":         67,
	}

	c.JSON(http.StatusOK, usage)
}
