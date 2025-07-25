
package api

import (
	"AdminiSoftware/internal/api/handlers"
	"AdminiSoftware/internal/api/middleware"
	"AdminiSoftware/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client, logger *utils.Logger) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, logger)

	// Public routes
	auth := router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/enable-2fa", authHandler.EnableTwoFactor)
		auth.POST("/verify-2fa", authHandler.VerifyTwoFactor)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Basic placeholder routes
		api.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "API is working"})
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})

	return router
}
