package main

import (
	"log"
	"net/http"

	"dailytrackr/gateway/proxy"
	"dailytrackr/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Add CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Gateway health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "dailytrackr-gateway",
			"status":  "healthy",
			"version": "1.0.0",
			"services": map[string]string{
				"user-service":         "http://localhost:" + cfg.UserServicePort,
				"activity-service":     "http://localhost:" + cfg.ActivityPort,
				"habit-service":        "http://localhost:" + cfg.HabitPort,
				"notification-service": "http://localhost:" + cfg.NotificationPort,
				"stat-service":         "http://localhost:" + cfg.StatPort,
				"ai-service":           "http://localhost:" + cfg.AIPort,
			},
		})
	})

	// Setup service proxies
	setupRoutes(r, cfg)

	// Start gateway
	port := ":" + cfg.GatewayPort
	log.Printf("🚀 DailyTrackr Gateway starting on port %s", port)
	log.Printf("📋 Available services:")
	log.Printf("   - User Service: http://localhost:%s", cfg.UserServicePort)
	log.Printf("   - Activity Service: http://localhost:%s", cfg.ActivityPort)
	log.Printf("   - Habit Service: http://localhost:%s", cfg.HabitPort)
	log.Printf("   - Stat Service: http://localhost:%s", cfg.StatPort)
	log.Printf("   - AI Service: http://localhost:%s", cfg.AIPort)
	log.Printf("🌐 Gateway URL: http://localhost:%s", cfg.GatewayPort)

	log.Fatal(r.Run(port))
}

// setupRoutes configures all microservice routes
func setupRoutes(r *gin.Engine, cfg *config.Config) {
	// User Service routes
	userProxy := proxy.NewServiceProxy("http://localhost:" + cfg.UserServicePort)

	userRoutes := r.Group("/api/users")
	{
		userRoutes.Any("/health", userProxy.ProxyRequest)
		userRoutes.Any("/auth/*path", userProxy.ProxyRequest)
		userRoutes.Any("/api/v1/users/*path", userProxy.ProxyRequest)
	}

	// Activity Service routes
	activityProxy := proxy.NewServiceProxy("http://localhost:" + cfg.ActivityPort)

	activityRoutes := r.Group("/api/activities")
	{
		activityRoutes.Any("/health", activityProxy.ProxyRequest)
		activityRoutes.Any("/api/v1/activities/*path", activityProxy.ProxyRequest)
	}

	// Habit Service routes
	habitProxy := proxy.NewServiceProxy("http://localhost:" + cfg.HabitPort)

	habitRoutes := r.Group("/api/habits")
	{
		habitRoutes.Any("/health", habitProxy.ProxyRequest)
		habitRoutes.Any("/api/v1/habits/*path", habitProxy.ProxyRequest)
	}

	// Statistics Service routes (NEW!)
	statProxy := proxy.NewServiceProxy("http://localhost:" + cfg.StatPort)

	statRoutes := r.Group("/api/stats")
	{
		statRoutes.Any("/health", statProxy.ProxyRequest)
		statRoutes.Any("/api/v1/stats/*path", statProxy.ProxyRequest)
	}

	// AI Service routes
	aiProxy := proxy.NewServiceProxy("http://localhost:" + cfg.AIPort)

	aiRoutes := r.Group("/api/ai")
	{
		aiRoutes.Any("/health", aiProxy.ProxyRequest)
		aiRoutes.Any("/api/v1/ai/*path", aiProxy.ProxyRequest)
	}
}
