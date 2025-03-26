package main

import (
	"log"
	"os"
	"r2-gallery/config"
	"r2-gallery/routes"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Initialize database connection
	config.InitDatabase()

	// Initialize R2 storage
	config.InitR2()

	// Set Gin mode based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = gin.DebugMode
	}
	gin.SetMode(ginMode)

	// Create Gin router
	router := gin.Default()

	// Setup CORS if needed
	router.Use(gin.Recovery()) // 避免崩溃
	router.Use(cors.Default()) // 允许跨域

	// API routes
	api := router.Group("/api")
	{
		// Setup auth routes
		routes.SetupAuthRoutes(api)

		// Setup image routes
		routes.SetupImageRoutes(api)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
