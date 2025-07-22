
package user

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WordPressHandler struct {
	db *gorm.DB
}

func NewWordPressHandler(db *gorm.DB) *WordPressHandler {
	return &WordPressHandler{db: db}
}

func (h *WordPressHandler) GetWordPressSites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var sites []models.WordPressSite
	if err := h.db.Where("user_id = ?", userID).Find(&sites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch WordPress sites"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sites": sites})
}

func (h *WordPressHandler) CreateWordPressSite(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		DomainID    uint   `json:"domain_id" binding:"required"`
		Path        string `json:"path"`
		SiteTitle   string `json:"site_title" binding:"required"`
		AdminUser   string `json:"admin_user" binding:"required"`
		AdminPass   string `json:"admin_pass" binding:"required"`
		AdminEmail  string `json:"admin_email" binding:"required"`
		DatabaseID  uint   `json:"database_id" binding:"required"`
		Version     string `json:"version"`
		Theme       string `json:"theme"`
		Plugins     []string `json:"plugins"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site := models.WordPressSite{
		UserID:     userID,
		DomainID:   req.DomainID,
		Path:       req.Path,
		SiteTitle:  req.SiteTitle,
		AdminUser:  req.AdminUser,
		AdminEmail: req.AdminEmail,
		DatabaseID: req.DatabaseID,
		Version:    req.Version,
		Theme:      req.Theme,
		Status:     "installing",
	}

	if err := h.db.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create WordPress site"})
		return
	}

	// Simulate WordPress installation
	go func() {
		// Install WordPress core files
		// Setup database
		// Configure WordPress
		site.Status = "active"
		h.db.Save(&site)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "WordPress installation started", "site": site})
}

func (h *WordPressHandler) UpdateWordPress(c *gin.Context) {
	userID := c.GetUint("user_id")
	siteID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var site models.WordPressSite
	if err := h.db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordPress site not found"})
		return
	}

	var req struct {
		Version string `json:"version"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	site.Status = "updating"
	site.Version = req.Version

	if err := h.db.Save(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start WordPress update"})
		return
	}

	// Simulate update process
	go func() {
		site.Status = "active"
		h.db.Save(&site)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "WordPress update started", "site": site})
}

func (h *WordPressHandler) GetWordPressPlugins(c *gin.Context) {
	siteID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")

	var site models.WordPressSite
	if err := h.db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordPress site not found"})
		return
	}

	var plugins []models.WordPressPlugin
	if err := h.db.Where("site_id = ?", siteID).Find(&plugins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch plugins"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plugins": plugins})
}

func (h *WordPressHandler) InstallPlugin(c *gin.Context) {
	userID := c.GetUint("user_id")
	siteID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var site models.WordPressSite
	if err := h.db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordPress site not found"})
		return
	}

	var req struct {
		Name        string `json:"name" binding:"required"`
		Version     string `json:"version"`
		Source      string `json:"source"`
		AutoActivate bool  `json:"auto_activate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plugin := models.WordPressPlugin{
		SiteID:  uint(siteID),
		Name:    req.Name,
		Version: req.Version,
		Status:  "installing",
	}

	if err := h.db.Create(&plugin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to install plugin"})
		return
	}

	// Simulate plugin installation
	go func() {
		if req.AutoActivate {
			plugin.Status = "active"
		} else {
			plugin.Status = "inactive"
		}
		h.db.Save(&plugin)
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Plugin installation started", "plugin": plugin})
}

func (h *WordPressHandler) GetWordPressThemes(c *gin.Context) {
	siteID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	userID := c.GetUint("user_id")

	var site models.WordPressSite
	if err := h.db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordPress site not found"})
		return
	}

	var themes []models.WordPressTheme
	if err := h.db.Where("site_id = ?", siteID).Find(&themes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch themes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"themes": themes})
}

func (h *WordPressHandler) ActivateTheme(c *gin.Context) {
	userID := c.GetUint("user_id")
	siteID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	themeID, _ := strconv.ParseUint(c.Param("theme_id"), 10, 32)

	var site models.WordPressSite
	if err := h.db.Where("id = ? AND user_id = ?", siteID, userID).First(&site).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "WordPress site not found"})
		return
	}

	// Deactivate current theme
	h.db.Model(&models.WordPressTheme{}).Where("site_id = ? AND status = 'active'", siteID).Update("status", "inactive")

	// Activate selected theme
	var theme models.WordPressTheme
	if err := h.db.Where("id = ? AND site_id = ?", themeID, siteID).First(&theme).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Theme not found"})
		return
	}

	theme.Status = "active"
	site.Theme = theme.Name

	if err := h.db.Save(&theme).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to activate theme"})
		return
	}

	if err := h.db.Save(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update site theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme activated successfully", "theme": theme})
}
