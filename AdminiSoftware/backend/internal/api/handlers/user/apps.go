
package user

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppsHandler struct {
	db *gorm.DB
}

func NewAppsHandler(db *gorm.DB) *AppsHandler {
	return &AppsHandler{db: db}
}

func (h *AppsHandler) GetApplications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var apps []models.Application
	if err := h.db.Where("user_id = ?", userID).Find(&apps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"applications": apps})
}

func (h *AppsHandler) GetAvailableApplications(c *gin.Context) {
	availableApps := []map[string]interface{}{
		{
			"name":        "WordPress",
			"description": "Popular blogging platform and CMS",
			"version":     "6.4",
			"category":    "Blog",
			"icon":        "/assets/wordpress.png",
		},
		{
			"name":        "Joomla",
			"description": "Flexible content management system",
			"version":     "4.4",
			"category":    "CMS",
			"icon":        "/assets/joomla.png",
		},
		{
			"name":        "Drupal",
			"description": "Open-source content management framework",
			"version":     "10.1",
			"category":    "CMS",
			"icon":        "/assets/drupal.png",
		},
		{
			"name":        "PrestaShop",
			"description": "E-commerce solution",
			"version":     "8.1",
			"category":    "E-commerce",
			"icon":        "/assets/prestashop.png",
		},
		{
			"name":        "Magento",
			"description": "Enterprise e-commerce platform",
			"version":     "2.4",
			"category":    "E-commerce",
			"icon":        "/assets/magento.png",
		},
	}

	c.JSON(http.StatusOK, gin.H{"applications": availableApps})
}

func (h *AppsHandler) InstallApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		Name        string `json:"name" binding:"required"`
		DomainID    uint   `json:"domain_id" binding:"required"`
		Path        string `json:"path"`
		DatabaseID  uint   `json:"database_id"`
		AdminUser   string `json:"admin_user"`
		AdminPass   string `json:"admin_pass"`
		AdminEmail  string `json:"admin_email"`
		SiteTitle   string `json:"site_title"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	app := models.Application{
		UserID:      userID,
		DomainID:    req.DomainID,
		Name:        req.Name,
		Path:        req.Path,
		DatabaseID:  req.DatabaseID,
		Status:      "installing",
		AdminUser:   req.AdminUser,
		AdminEmail:  req.AdminEmail,
		SiteTitle:   req.SiteTitle,
		Description: req.Description,
	}

	if err := h.db.Create(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initiate application installation"})
		return
	}

	// Simulate installation process (in real implementation, this would be asynchronous)
	go func() {
		// Simulate installation time
		// time.Sleep(30 * time.Second)
		app.Status = "installed"
		h.db.Save(&app)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Application installation started", "application": app})
}

func (h *AppsHandler) UpdateApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var app models.Application
	if err := h.db.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	app.Status = "updating"
	if err := h.db.Save(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start application update"})
		return
	}

	// Simulate update process
	go func() {
		// time.Sleep(15 * time.Second)
		app.Status = "installed"
		h.db.Save(&app)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Application update started", "application": app})
}

func (h *AppsHandler) UninstallApplication(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var app models.Application
	if err := h.db.Where("id = ? AND user_id = ?", appID, userID).First(&app).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	if err := h.db.Delete(&app).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to uninstall application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application uninstalled successfully"})
}
