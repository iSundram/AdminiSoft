
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
	logger.Info("Starting AdminiSoftware server...")
	
	// Initialize database
	db, err := config.InitDatabase(cfg)
	if err != nil {
		logger.Error("Failed to initialize database: " + err.Error())
		log.Fatal("Failed to initialize database:", err)
	}
	logger.Info("Database initialized successfully")
	
	// Initialize Redis
	redis := config.InitRedis(cfg)
	if redis == nil {
		logger.Error("Failed to initialize Redis connection")
		log.Fatal("Failed to initialize Redis connection")
	}
	logger.Info("Redis initialized successfully")
	
	// Setup API routes
	router := api.SetupRouter(db, redis, logger)
	logger.Info("API routes initialized successfully")
	
	// Start server
	serverAddr := "0.0.0.0:" + cfg.Port
	logger.Info("Starting AdminiSoftware server on " + serverAddr)
	logger.Info("Environment: " + cfg.Environment)
	
	if err := router.Run(serverAddr); err != nil {
		logger.Error("Failed to start server: " + err.Error())
		log.Fatal("Failed to start server:", err)
	}
}
