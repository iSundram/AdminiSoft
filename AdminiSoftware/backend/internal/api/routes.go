
package api

import (
	"AdminiSoftware/internal/api/handlers"
	"AdminiSoftware/internal/api/middleware"
	"AdminiSoftware/internal/auth"
	"AdminiSoftware/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client, logger *utils.Logger) *gin.Engine {
	router := gin.New()

	// Initialize managers
	jwtManager := auth.NewJWTManager("your-secret-key-change-this")
	bruteForce := auth.NewBruteForceProtection(redis)
	rateLimiter := middleware.NewRateLimiter(redis)

	// Global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecureHeaders())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.ErrorLogger(logger))
	router.Use(gin.Recovery())

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, jwtManager, bruteForce, logger)

	// Public routes
	auth := router.Group("/api/auth")
	auth.Use(rateLimiter.AuthRateLimit())
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
		auth.POST("/reset-password", authHandler.ResetPassword)
		auth.POST("/verify-2fa", authHandler.Verify2FA)
	}

	// Protected routes
	api := router.Group("/api")
	api.Use(rateLimiter.APIRateLimit())
	api.Use(middleware.AuthMiddleware(jwtManager))
	{
		// User routes
		user := api.Group("/user")
		{
			user.GET("/profile", authHandler.GetProfile)
			user.PUT("/profile", authHandler.UpdateProfile)
			user.POST("/change-password", authHandler.ChangePassword)
			user.POST("/enable-2fa", authHandler.Enable2FA)
			user.POST("/disable-2fa", authHandler.Disable2FA)
		}

		// Admin routes
		admin := api.Group("/admin")
		admin.Use(middleware.RequireRole("admin"))
		{
			// Account management
			accounts := admin.Group("/accounts")
			{
				accounts.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List accounts"}) })
				accounts.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create account"}) })
				accounts.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update account"}) })
				accounts.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete account"}) })
			}

			// Package management
			packages := admin.Group("/packages")
			{
				packages.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List packages"}) })
				packages.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create package"}) })
				packages.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update package"}) })
				packages.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete package"}) })
			}

			// System management
			system := admin.Group("/system")
			{
				system.GET("/stats", func(c *gin.Context) { c.JSON(200, gin.H{"message": "System stats"}) })
				system.GET("/services", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List services"}) })
				system.POST("/services/:name/restart", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Restart service"}) })
			}
		}

		// Reseller routes
		reseller := api.Group("/reseller")
		reseller.Use(middleware.RequireRole("reseller", "admin"))
		{
			// Reseller account management
			accounts := reseller.Group("/accounts")
			{
				accounts.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List reseller accounts"}) })
				accounts.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create reseller account"}) })
			}
		}

		// User panel routes
		panel := api.Group("/panel")
		{
			// Domain management
			domains := panel.Group("/domains")
			{
				domains.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List domains"}) })
				domains.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create domain"}) })
				domains.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update domain"}) })
				domains.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete domain"}) })
			}

			// File management
			files := panel.Group("/files")
			{
				files.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List files"}) })
				files.POST("/upload", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Upload file"}) })
				files.DELETE("/:path", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete file"}) })
			}

			// Email management
			emails := panel.Group("/emails")
			{
				emails.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List emails"}) })
				emails.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create email"}) })
				emails.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update email"}) })
				emails.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete email"}) })
			}

			// Database management
			databases := panel.Group("/databases")
			{
				databases.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List databases"}) })
				databases.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create database"}) })
				databases.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update database"}) })
				databases.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete database"}) })
			}

			// Application management
			apps := panel.Group("/applications")
			{
				apps.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List applications"}) })
				apps.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Install application"}) })
				apps.PUT("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Update application"}) })
				apps.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Uninstall application"}) })
			}

			// SSL management
			ssl := panel.Group("/ssl")
			{
				ssl.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List SSL certificates"}) })
				ssl.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Install SSL certificate"}) })
				ssl.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Remove SSL certificate"}) })
			}

			// Backup management
			backups := panel.Group("/backups")
			{
				backups.GET("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "List backups"}) })
				backups.POST("/", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Create backup"}) })
				backups.POST("/:id/restore", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Restore backup"}) })
				backups.DELETE("/:id", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Delete backup"}) })
			}

			// Statistics
			stats := panel.Group("/stats")
			{
				stats.GET("/usage", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Usage statistics"}) })
				stats.GET("/bandwidth", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Bandwidth statistics"}) })
				stats.GET("/visitors", func(c *gin.Context) { c.JSON(200, gin.H{"message": "Visitor statistics"}) })
			}
		}
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
