
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
