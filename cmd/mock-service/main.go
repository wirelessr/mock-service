package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"mock-service/internal/config"
	"mock-service/internal/handler"
	"mock-service/internal/logger"
	"mock-service/internal/matcher"
	"mock-service/internal/response"

	"github.com/gin-gonic/gin"
)

func main() {
	// Parse command line flags
	var configFile string
	var port string

	flag.StringVar(&configFile, "config", "config.json", "Path to configuration file")
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.Parse()

	// Initialize components
	configManager := config.NewConfigManager()
	pathMatcher := matcher.NewPathMatcher()
	responseBuilder := response.NewResponseBuilder()
	appLogger := logger.NewLogger()

	// Load configuration
	if err := configManager.LoadConfig(configFile); err != nil {
		log.Fatalf("Failed to load configuration from %s: %v", configFile, err)
	}

	// Create universal handler
	universalHandler := handler.NewUniversalHandler(
		configManager,
		pathMatcher,
		responseBuilder,
		appLogger,
	)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode) // Disable Gin debug output
	router := gin.New()

	// Use custom middleware for logging (we handle logging in our handler)
	router.Use(gin.Recovery())

	// Register health check endpoint first (before wildcard routes)
	router.GET("/health", func(c *gin.Context) {
		const healthStatusCode = 200
		c.JSON(healthStatusCode, gin.H{"status": "healthy", "service": "mock-service"})
	})

	// Register universal handler for all other paths and methods
	// Note: NoRoute handles requests that don't match any registered routes
	router.NoRoute(universalHandler.HandleRequest)

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Mock service starting on port %s\n", port)
		fmt.Printf("Configuration loaded from: %s\n", configFile)

		if err := router.Run(":" + port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	fmt.Println("\nShutting down mock service...")
}
