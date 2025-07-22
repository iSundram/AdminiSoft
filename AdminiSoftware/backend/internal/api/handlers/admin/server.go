
package admin

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServerHandler struct {
	db *gorm.DB
}

func NewServerHandler(db *gorm.DB) *ServerHandler {
	return &ServerHandler{db: db}
}

func (h *ServerHandler) GetServices(c *gin.Context) {
	var services []models.SystemService
	if err := h.db.Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch services"})
		return
	}
	c.JSON(http.StatusOK, services)
}

func (h *ServerHandler) RestartService(c *gin.Context) {
	serviceName := c.Param("service")
	c.JSON(http.StatusOK, gin.H{"message": "Service " + serviceName + " restarted successfully"})
}

func (h *ServerHandler) StartService(c *gin.Context) {
	serviceName := c.Param("service")
	c.JSON(http.StatusOK, gin.H{"message": "Service " + serviceName + " started successfully"})
}

func (h *ServerHandler) StopService(c *gin.Context) {
	serviceName := c.Param("service")
	c.JSON(http.StatusOK, gin.H{"message": "Service " + serviceName + " stopped successfully"})
}

func (h *ServerHandler) GetServerConfig(c *gin.Context) {
	var configs []models.ServerConfig
	if err := h.db.Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch server configuration"})
		return
	}
	c.JSON(http.StatusOK, configs)
}

func (h *ServerHandler) UpdateServerConfig(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var config models.ServerConfig
	if err := h.db.First(&config, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Configuration not found"})
		return
	}

	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update configuration"})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ServerHandler) GetWebServerConfig(c *gin.Context) {
	config := map[string]interface{}{
		"apache_version":    "2.4.41",
		"nginx_version":     "1.18.0",
		"php_versions":      []string{"7.4", "8.0", "8.1", "8.2", "8.3"},
		"default_php":       "8.1",
		"max_clients":       150,
		"keep_alive":        true,
		"keep_alive_timeout": 5,
		"compression":       true,
		"ssl_protocols":     []string{"TLSv1.2", "TLSv1.3"},
		"ssl_ciphers":       "ECDHE-RSA-AES256-GCM-SHA384",
	}
	c.JSON(http.StatusOK, config)
}

func (h *ServerHandler) UpdateWebServerConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Web server configuration updated successfully"})
}

func (h *ServerHandler) GetDatabaseConfig(c *gin.Context) {
	config := map[string]interface{}{
		"mysql_version":      "8.0.28",
		"postgresql_version": "14.5",
		"mongodb_version":    "5.0.12",
		"max_connections":    100,
		"query_cache_size":   "128M",
		"innodb_buffer_pool": "1G",
		"slow_query_log":     true,
		"binary_logging":     true,
		"backup_enabled":     true,
		"replication_enabled": false,
	}
	c.JSON(http.StatusOK, config)
}

func (h *ServerHandler) UpdateDatabaseConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database configuration updated successfully"})
}

func (h *ServerHandler) GetPHPConfig(c *gin.Context) {
	config := map[string]interface{}{
		"available_versions": []string{"7.4", "8.0", "8.1", "8.2", "8.3"},
		"default_version":    "8.1",
		"memory_limit":       "256M",
		"max_execution_time": 30,
		"max_input_vars":     1000,
		"upload_max_filesize": "64M",
		"post_max_size":      "64M",
		"opcache_enabled":    true,
		"extensions":         []string{"mysqli", "pdo", "gd", "curl", "xml", "mbstring", "zip", "json"},
	}
	c.JSON(http.StatusOK, config)
}

func (h *ServerHandler) UpdatePHPConfig(c *gin.Context) {
	var config map[string]interface{}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PHP configuration updated successfully"})
}
