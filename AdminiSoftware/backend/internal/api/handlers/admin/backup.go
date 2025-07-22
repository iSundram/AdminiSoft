
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BackupHandler struct {
	db *gorm.DB
}

func NewBackupHandler(db *gorm.DB) *BackupHandler {
	return &BackupHandler{db: db}
}

func (h *BackupHandler) GetBackups(c *gin.Context) {
	var backups []models.Backup
	if err := h.db.Preload("User").Find(&backups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch backups"})
		return
	}
	c.JSON(http.StatusOK, backups)
}

func (h *BackupHandler) CreateBackup(c *gin.Context) {
	var backup models.Backup
	if err := c.ShouldBindJSON(&backup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	backup.Status = "pending"
	if err := h.db.Create(&backup).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup"})
		return
	}

	c.JSON(http.StatusCreated, backup)
}

func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var backup models.Backup
	if err := h.db.First(&backup, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Backup restoration initiated"})
}

func (h *BackupHandler) GetSchedules(c *gin.Context) {
	var schedules []models.BackupSchedule
	if err := h.db.Preload("User").Find(&schedules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch backup schedules"})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

func (h *BackupHandler) CreateSchedule(c *gin.Context) {
	var schedule models.BackupSchedule
	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup schedule"})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

func (h *BackupHandler) UpdateSchedule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var schedule models.BackupSchedule
	if err := h.db.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup schedule not found"})
		return
	}

	if err := c.ShouldBindJSON(&schedule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update backup schedule"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func (h *BackupHandler) DeleteSchedule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.BackupSchedule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete backup schedule"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Backup schedule deleted successfully"})
}
