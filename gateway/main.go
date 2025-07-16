package main

import (
	"log"
	"net/http"
	"time"

	"dailytrackr/gateway/proxy"
	"dailytrackr/shared/config"

	"github.com/gin-contrib/cors"
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

	// ================================================
	// ENHANCED CORS MIDDLEWARE - Critical Fix
	// ================================================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3000", "*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Additional manual CORS handling for edge cases
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept, X-Requested-With")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Gateway health check
	r.GET("/", func(c *gin.Context) {
		services := map[string]interface{}{
			"user-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.UserServicePort,
				"status": checkServiceHealth("http://localhost:" + cfg.UserServicePort + "/health"),
			},
			"activity-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.ActivityPort,
				"status": checkServiceHealth("http://localhost:" + cfg.ActivityPort + "/health"),
			},
			"habit-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.HabitPort,
				"status": checkServiceHealth("http://localhost:" + cfg.HabitPort + "/health"),
			},
			"stat-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.StatPort,
				"status": checkServiceHealth("http://localhost:" + cfg.StatPort + "/health"),
			},
			"ai-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.AIPort,
				"status": checkServiceHealth("http://localhost:" + cfg.AIPort + "/health"),
			},
		}

		// Count healthy services
		healthyCount := 0
		totalCount := len(services)
		for _, service := range services {
			if serviceMap, ok := service.(map[string]interface{}); ok {
				if status, exists := serviceMap["status"]; exists && status == "healthy" {
					healthyCount++
				}
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"service":          "dailytrackr-gateway",
			"status":           "healthy",
			"version":          "1.0.0",
			"healthy_services": healthyCount,
			"total_services":   totalCount,
			"services":         services,
			"message":          "DailyTrackr Gateway is running successfully! 🚀",
		})
	})

	// Setup service proxies with FIXED routing
	setupRoutesFixed(r, cfg)

	// Start gateway
	port := ":" + cfg.GatewayPort
	log.Printf("🚀 DailyTrackr Gateway starting on port %s", port)
	log.Printf("📋 Service endpoints:")
	log.Printf("   - User Service:         http://localhost:%s", cfg.UserServicePort)
	log.Printf("   - Activity Service:     http://localhost:%s", cfg.ActivityPort)
	log.Printf("   - Habit Service:        http://localhost:%s", cfg.HabitPort)
	log.Printf("   - Statistics Service:   http://localhost:%s", cfg.StatPort)
	log.Printf("   - AI Service:           http://localhost:%s", cfg.AIPort)
	log.Printf("🌐 Gateway URL: http://localhost:%s", cfg.GatewayPort)
	log.Printf("🔗 API Routes:")
	log.Printf("   - /api/users/*          -> User Service")
	log.Printf("   - /api/activities/*     -> Activity Service")
	log.Printf("   - /api/habits/*         -> Habit Service")
	log.Printf("   - /api/stats/*          -> Statistics Service")
	log.Printf("   - /api/ai/*             -> AI Service")

	log.Fatal(r.Run(port))
}

// setupRoutesFixed configures all microservice routes with proper mapping
func setupRoutesFixed(r *gin.Engine, cfg *config.Config) {
	// User Service routes
	userProxy := proxy.NewServiceProxy("http://localhost:" + cfg.UserServicePort)
	userRoutes := r.Group("/api/users")
	{
		userRoutes.Any("/health", userProxy.ProxyRequest)
		userRoutes.Any("/auth/*path", userProxy.ProxyRequest)
		userRoutes.Any("/api/v1/users/*path", userProxy.ProxyRequest)
		// FIXED: Direct auth routes without /api/users prefix
		r.Any("/auth/*path", userProxy.ProxyRequest)
	}

	// Activity Service routes
	activityProxy := proxy.NewServiceProxy("http://localhost:" + cfg.ActivityPort)
	activityRoutes := r.Group("/api/activities")
	{
		activityRoutes.Any("/health", activityProxy.ProxyRequest)
		activityRoutes.Any("/api/v1/activities", activityProxy.ProxyRequest)       // Exact match
		activityRoutes.Any("/api/v1/activities/*path", activityProxy.ProxyRequest) // With path
	}

	// Habit Service routes - ENHANCED ROUTING
	habitProxy := proxy.NewServiceProxy("http://localhost:" + cfg.HabitPort)
	habitRoutes := r.Group("/api/habits")
	{
		habitRoutes.Any("/health", habitProxy.ProxyRequest)
		habitRoutes.Any("/api/v1/habits", habitProxy.ProxyRequest)       // Exact match for GET /api/v1/habits
		habitRoutes.Any("/api/v1/habits/*path", habitProxy.ProxyRequest) // All other habit routes
		habitRoutes.Any("/api/v1/habit-logs/*path", habitProxy.ProxyRequest)
	}

	// Statistics Service routes
	statProxy := proxy.NewServiceProxy("http://localhost:" + cfg.StatPort)
	statRoutes := r.Group("/api/stats")
	{
		statRoutes.Any("/health", statProxy.ProxyRequest)
		statRoutes.Any("/api/v1/stats", statProxy.ProxyRequest)       // Exact match
		statRoutes.Any("/api/v1/stats/*path", statProxy.ProxyRequest) // With path
	}

	// AI Service routes
	aiProxy := proxy.NewServiceProxy("http://localhost:" + cfg.AIPort)
	aiRoutes := r.Group("/api/ai")
	{
		aiRoutes.Any("/health", aiProxy.ProxyRequest)
		aiRoutes.Any("/api/v1/ai", aiProxy.ProxyRequest)       // Exact match
		aiRoutes.Any("/api/v1/ai/*path", aiProxy.ProxyRequest) // With path
	}

	// Notification Service routes (for future implementation)
	notificationProxy := proxy.NewServiceProxy("http://localhost:" + cfg.NotificationPort)
	notificationRoutes := r.Group("/api/notifications")
	{
		notificationRoutes.Any("/health", notificationProxy.ProxyRequest)
		notificationRoutes.Any("/api/v1/notifications", notificationProxy.ProxyRequest)       // Exact match
		notificationRoutes.Any("/api/v1/notifications/*path", notificationProxy.ProxyRequest) // With path
	}

	// Add debugging route to test routing
	r.GET("/debug/routes", func(c *gin.Context) {
		routes := []map[string]string{
			{"pattern": "/api/users/*", "target": "User Service :" + cfg.UserServicePort},
			{"pattern": "/api/activities/*", "target": "Activity Service :" + cfg.ActivityPort},
			{"pattern": "/api/habits/*", "target": "Habit Service :" + cfg.HabitPort},
			{"pattern": "/api/stats/*", "target": "Statistics Service :" + cfg.StatPort},
			{"pattern": "/api/ai/*", "target": "AI Service :" + cfg.AIPort},
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "DailyTrackr Gateway Route Mapping",
			"routes":  routes,
			"note":    "All routes proxy to their respective microservices",
		})
	})
}

// checkServiceHealth checks if a service is healthy
func checkServiceHealth(url string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "down"
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "healthy"
	}

	return "unhealthy"
}
