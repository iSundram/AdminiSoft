
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DNSHandler struct {
	db *gorm.DB
}

func NewDNSHandler(db *gorm.DB) *DNSHandler {
	return &DNSHandler{db: db}
}

func (h *DNSHandler) GetZones(c *gin.Context) {
	var zones []models.DNSZone
	if err := h.db.Preload("Records").Find(&zones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DNS zones"})
		return
	}
	c.JSON(http.StatusOK, zones)
}

func (h *DNSHandler) CreateZone(c *gin.Context) {
	var zone models.DNSZone
	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&zone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DNS zone"})
		return
	}

	c.JSON(http.StatusCreated, zone)
}

func (h *DNSHandler) UpdateZone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var zone models.DNSZone
	if err := h.db.First(&zone, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS zone not found"})
		return
	}

	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&zone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DNS zone"})
		return
	}

	c.JSON(http.StatusOK, zone)
}

func (h *DNSHandler) DeleteZone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.DNSZone{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS zone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DNS zone deleted successfully"})
}

func (h *DNSHandler) CreateRecord(c *gin.Context) {
	var record models.DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DNS record"})
		return
	}

	c.JSON(http.StatusCreated, record)
}

func (h *DNSHandler) UpdateRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.DNSRecord
	if err := h.db.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS record not found"})
		return
	}

	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DNS record"})
		return
	}

	c.JSON(http.StatusOK, record)
}

func (h *DNSHandler) DeleteRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.DNSRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS record"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DNS record deleted successfully"})
}
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DNSHandler struct {
	db *gorm.DB
}

func NewDNSHandler(db *gorm.DB) *DNSHandler {
	return &DNSHandler{db: db}
}

func (h *DNSHandler) GetDNSZones(c *gin.Context) {
	var zones []models.DNSZone
	if err := h.db.Find(&zones).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DNS zones"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dns_zones": zones})
}

func (h *DNSHandler) CreateDNSZone(c *gin.Context) {
	var zone models.DNSZone
	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&zone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DNS zone"})
		return
	}

	c.JSON(http.StatusCreated, zone)
}

func (h *DNSHandler) UpdateDNSZone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var zone models.DNSZone
	if err := h.db.First(&zone, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DNS zone not found"})
		return
	}

	if err := c.ShouldBindJSON(&zone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&zone).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update DNS zone"})
		return
	}

	c.JSON(http.StatusOK, zone)
}

func (h *DNSHandler) DeleteDNSZone(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&models.DNSZone{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete DNS zone"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "DNS zone deleted successfully"})
}

func (h *DNSHandler) GetDNSRecords(c *gin.Context) {
	zoneID := c.Param("zone_id")
	var records []models.DNSRecord
	if err := h.db.Where("zone_id = ?", zoneID).Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch DNS records"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"dns_records": records})
}

func (h *DNSHandler) CreateDNSRecord(c *gin.Context) {
	var record models.DNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DNS record"})
		return
	}

	c.JSON(http.StatusCreated, record)
}
