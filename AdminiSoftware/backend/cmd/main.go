
package main

import (
	"AdminiSoftware/internal/config"
	"AdminiSoftware/internal/api"
	"AdminiSoftware/internal/utils"
	"log"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Initialize logger
	logger := utils.NewLogger()
	
	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	
	// Initialize Redis
	redis := config.InitRedis(cfg)
	
	// Setup API routes
	router := api.SetupRouter(db, redis, logger)
	
	// Start server
	logger.Info("Starting AdminiSoftware server on port " + cfg.Port)
	if err := router.Run("0.0.0.0:" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
package main

import (
	"AdminiSoftware/internal/api"
	"AdminiSoftware/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db := config.InitDatabase(cfg)

	// Initialize Redis
	redis := config.InitRedis(cfg)

	// Setup Gin router
	r := gin.Default()

	// Setup API routes
	api.SetupRoutes(r, db, redis)

	// Create server
	server := &http.Server{
		Addr:           "0.0.0.0:5000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server starting on port 5000...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
