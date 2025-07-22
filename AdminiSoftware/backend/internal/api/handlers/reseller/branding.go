
package reseller

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BrandingHandler struct {
	db *gorm.DB
}

func NewBrandingHandler(db *gorm.DB) *BrandingHandler {
	return &BrandingHandler{db: db}
}

func (h *BrandingHandler) GetBranding(c *gin.Context) {
	userID := c.GetUint("user_id")
	var branding models.Branding
	if err := h.db.Where("user_id = ?", userID).First(&branding).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default branding
			branding = models.Branding{
				UserID:        userID,
				CompanyName:   "AdminiSoftware",
				LogoURL:       "/assets/logos/adminisoftware-logo.svg",
				ThemeColor:    "#3B82F6",
				ShowPoweredBy: true,
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branding"})
			return
		}
	}
	c.JSON(http.StatusOK, branding)
}

func (h *BrandingHandler) UpdateBranding(c *gin.Context) {
	userID := c.GetUint("user_id")
	var branding models.Branding
	
	// Try to find existing branding
	err := h.db.Where("user_id = ?", userID).First(&branding).Error
	if err == gorm.ErrRecordNotFound {
		// Create new branding
		branding.UserID = userID
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branding"})
		return
	}

	if err := c.ShouldBindJSON(&branding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branding.UserID = userID
	if err := h.db.Save(&branding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branding"})
		return
	}

	c.JSON(http.StatusOK, branding)
}

func (h *BrandingHandler) GetThemes(c *gin.Context) {
	themes := []map[string]interface{}{
		{
			"id":          "default",
			"name":        "AdminiSoftware Default",
			"description": "Clean and modern default theme",
			"preview":     "/assets/themes/default-preview.png",
			"colors": map[string]string{
				"primary":   "#3B82F6",
				"secondary": "#6B7280",
				"accent":    "#10B981",
				"background": "#F9FAFB",
			},
		},
		{
			"id":          "cpanel",
			"name":        "cPanel Style",
			"description": "Classic cPanel look and feel",
			"preview":     "/assets/themes/cpanel-preview.png",
			"colors": map[string]string{
				"primary":   "#2563EB",
				"secondary": "#475569",
				"accent":    "#059669",
				"background": "#FFFFFF",
			},
		},
		{
			"id":          "whm",
			"name":        "WHM Style",
			"description": "Professional WHM-inspired theme",
			"preview":     "/assets/themes/whm-preview.png",
			"colors": map[string]string{
				"primary":   "#1F2937",
				"secondary": "#6B7280",
				"accent":    "#EF4444",
				"background": "#F3F4F6",
			},
		},
		{
			"id":          "dark",
			"name":        "Dark Theme",
			"description": "Easy on the eyes dark theme",
			"preview":     "/assets/themes/dark-preview.png",
			"colors": map[string]string{
				"primary":   "#3B82F6",
				"secondary": "#9CA3AF",
				"accent":    "#F59E0B",
				"background": "#111827",
			},
		},
	}
	c.JSON(http.StatusOK, themes)
}

func (h *BrandingHandler) UpdateTheme(c *gin.Context) {
	userID := c.GetUint("user_id")
	var request struct {
		ThemeID string `json:"theme_id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var branding models.Branding
	err := h.db.Where("user_id = ?", userID).First(&branding).Error
	if err == gorm.ErrRecordNotFound {
		branding.UserID = userID
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branding"})
		return
	}

	// Update theme-specific settings based on theme ID
	switch request.ThemeID {
	case "cpanel":
		branding.ThemeColor = "#2563EB"
		branding.CustomCSS = ".cpanel-theme { /* cPanel specific styles */ }"
	case "whm":
		branding.ThemeColor = "#1F2937"
		branding.CustomCSS = ".whm-theme { /* WHM specific styles */ }"
	case "dark":
		branding.ThemeColor = "#3B82F6"
		branding.CustomCSS = ".dark-theme { /* Dark theme styles */ }"
	default:
		branding.ThemeColor = "#3B82F6"
		branding.CustomCSS = ""
	}

	if err := h.db.Save(&branding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update theme"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Theme updated successfully"})
}

func (h *BrandingHandler) UploadLogo(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, header, err := c.Request.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get uploaded file"})
		return
	}
	defer file.Close()

	// In a real implementation, you would save the file to storage
	// For now, we'll just simulate saving and return a URL
	logoURL := "/uploads/logos/" + strconv.Itoa(int(userID)) + "_" + header.Filename

	var branding models.Branding
	err = h.db.Where("user_id = ?", userID).First(&branding).Error
	if err == gorm.ErrRecordNotFound {
		branding.UserID = userID
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branding"})
		return
	}

	branding.LogoURL = logoURL
	if err := h.db.Save(&branding).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save logo URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Logo uploaded successfully",
		"logo_url": logoURL,
	})
}

func (h *BrandingHandler) GetCustomLinks(c *gin.Context) {
	userID := c.GetUint("user_id")
	// In a real implementation, you would fetch from database
	customLinks := []map[string]interface{}{
		{
			"id":     1,
			"title":  "Support Portal",
			"url":    "https://support.example.com",
			"target": "_blank",
			"icon":   "support",
			"order":  1,
		},
		{
			"id":     2,
			"title":  "Knowledge Base",
			"url":    "https://kb.example.com",
			"target": "_blank",
			"icon":   "book",
			"order":  2,
		},
	}
	c.JSON(http.StatusOK, customLinks)
}

func (h *BrandingHandler) UpdateCustomLinks(c *gin.Context) {
	userID := c.GetUint("user_id")
	var links []map[string]interface{}
	if err := c.ShouldBindJSON(&links); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real implementation, you would save to database
	c.JSON(http.StatusOK, gin.H{
		"message": "Custom links updated successfully",
		"user_id": userID,
		"links":   links,
	})
}
