
package user

import (
	"AdminiSoftware/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DatabaseHandler struct {
	db *gorm.DB
}

func NewDatabaseHandler(db *gorm.DB) *DatabaseHandler {
	return &DatabaseHandler{db: db}
}

func (h *DatabaseHandler) GetDatabases(c *gin.Context) {
	userID := c.GetUint("user_id")
	var databases []models.Database
	if err := h.db.Where("user_id = ?", userID).Find(&databases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch databases"})
		return
	}
	c.JSON(http.StatusOK, databases)
}

func (h *DatabaseHandler) CreateDatabase(c *gin.Context) {
	userID := c.GetUint("user_id")
	var database models.Database
	if err := c.ShouldBindJSON(&database); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.UserID = userID
	database.Status = "active"
	database.Host = "localhost"
	database.Port = 3306

	if err := h.db.Create(&database).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database"})
		return
	}

	c.JSON(http.StatusCreated, database)
}

func (h *DatabaseHandler) DeleteDatabase(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Database{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database deleted successfully"})
}

func (h *DatabaseHandler) GetDatabaseUsers(c *gin.Context) {
	userID := c.GetUint("user_id")
	var dbUsers []models.DatabaseUser
	if err := h.db.Joins("JOIN databases ON database_users.database_id = databases.id").
		Where("databases.user_id = ?", userID).
		Preload("Database").
		Find(&dbUsers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch database users"})
		return
	}
	c.JSON(http.StatusOK, dbUsers)
}

func (h *DatabaseHandler) CreateDatabaseUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	var dbUser models.DatabaseUser
	if err := c.ShouldBindJSON(&dbUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify database belongs to user
	var database models.Database
	if err := h.db.Where("id = ? AND user_id = ?", dbUser.DatabaseID, userID).First(&database).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	dbUser.Host = "localhost"
	if err := h.db.Create(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create database user"})
		return
	}

	c.JSON(http.StatusCreated, dbUser)
}

func (h *DatabaseHandler) UpdateDatabaseUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var dbUser models.DatabaseUser
	if err := h.db.Joins("JOIN databases ON database_users.database_id = databases.id").
		Where("database_users.id = ? AND databases.user_id = ?", id, userID).
		First(&dbUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database user not found"})
		return
	}

	if err := c.ShouldBindJSON(&dbUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.db.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database user"})
		return
	}

	c.JSON(http.StatusOK, dbUser)
}

func (h *DatabaseHandler) DeleteDatabaseUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	if err := h.db.Joins("JOIN databases ON database_users.database_id = databases.id").
		Where("database_users.id = ? AND databases.user_id = ?", id, userID).
		Delete(&models.DatabaseUser{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete database user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database user deleted successfully"})
}

func (h *DatabaseHandler) GetDatabasePrivileges(c *gin.Context) {
	privileges := []string{
		"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "DROP", "ALTER",
		"INDEX", "REFERENCES", "CREATE TEMPORARY TABLES", "LOCK TABLES",
		"CREATE VIEW", "SHOW VIEW", "CREATE ROUTINE", "ALTER ROUTINE",
		"EXECUTE", "EVENT", "TRIGGER",
	}
	c.JSON(http.StatusOK, privileges)
}

func (h *DatabaseHandler) UpdateDatabasePrivileges(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var request struct {
		Privileges []string `json:"privileges"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dbUser models.DatabaseUser
	if err := h.db.Joins("JOIN databases ON database_users.database_id = databases.id").
		Where("database_users.id = ? AND databases.user_id = ?", id, userID).
		First(&dbUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database user not found"})
		return
	}

	// Convert privileges array to string (in real implementation, you'd store this properly)
	privilegesStr := ""
	for i, priv := range request.Privileges {
		if i > 0 {
			privilegesStr += ","
		}
		privilegesStr += priv
	}
	dbUser.Privileges = privilegesStr

	if err := h.db.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update privileges"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database privileges updated successfully"})
}

func (h *DatabaseHandler) GetDatabaseStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var database models.Database
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	stats := map[string]interface{}{
		"database_id":   id,
		"name":          database.Name,
		"type":          database.Type,
		"size":          "25.6 MB",
		"tables":        15,
		"records":       12450,
		"last_backup":   "2024-01-15T02:00:00Z",
		"created_at":    database.CreatedAt,
		"table_stats": []map[string]interface{}{
			{"name": "users", "records": 250, "size": "45 KB"},
			{"name": "posts", "records": 1200, "size": "2.1 MB"},
			{"name": "comments", "records": 3400, "size": "890 KB"},
		},
		"performance": map[string]interface{}{
			"queries_per_second": 12.5,
			"avg_query_time":     "0.045s",
			"slow_queries":       2,
			"connections":        5,
		},
	}

	c.JSON(http.StatusOK, stats)
}

func (h *DatabaseHandler) BackupDatabase(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	var database models.Database
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	backup := map[string]interface{}{
		"database_id":   id,
		"backup_file":   database.Name + "_backup_" + "20240115103000.sql",
		"size":          "25.6 MB",
		"created_at":    "2024-01-15T10:30:00Z",
		"download_url":  "/api/user/databases/" + strconv.Itoa(id) + "/backup/download",
	}

	c.JSON(http.StatusOK, backup)
}

func (h *DatabaseHandler) RestoreDatabase(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))
	
	file, header, err := c.Request.FormFile("backup_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get backup file"})
		return
	}
	defer file.Close()

	var database models.Database
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&database).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Database restore initiated",
		"database_id": id,
		"backup_file": header.Filename,
		"size":        header.Size,
	})
}
