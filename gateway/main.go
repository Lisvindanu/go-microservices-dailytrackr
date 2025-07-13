package main

import (
	"log"
	"net/http"
	"time"

	"dailytrackr/gateway/proxy"
	"dailytrackr/shared/config"

	"github.com/gin-contrib/cors" // <-- 1. IMPORT PUSTAKA CORS
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

	// 2. GUNAKAN MIDDLEWARE CORS YANG LEBIH ANDAL
	// Ini akan menangani preflight request (OPTIONS) dengan benar.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Untuk development, bisa diganti dengan domain frontend nanti
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Gateway health check with service status
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
			"notification-service": map[string]interface{}{
				"url":    "http://localhost:" + cfg.NotificationPort,
				"status": checkServiceHealth("http://localhost:" + cfg.NotificationPort + "/health"),
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
			"endpoints": map[string]string{
				"users":         "/api/users/*",
				"activities":    "/api/activities/*",
				"habits":        "/api/habits/*",
				"statistics":    "/api/stats/*",
				"ai":            "/api/ai/*",
				"notifications": "/api/notifications/*",
			},
			"documentation": "/api/docs",
		})
	})

	// Setup service proxies
	setupRoutes(r, cfg)

	// Add utility routes
	setupUtilityRoutes(r, cfg)

	// Start gateway
	port := ":" + cfg.GatewayPort
	log.Printf("噫 DailyTrackr Gateway starting on port %s", port)
	log.Printf("搭 Available services:")
	log.Printf("   - User Service:         http://localhost:%s", cfg.UserServicePort)
	log.Printf("   - Activity Service:     http://localhost:%s", cfg.ActivityPort)
	log.Printf("   - Habit Service:        http://localhost:%s", cfg.HabitPort)
	log.Printf("   - Statistics Service:   http://localhost:%s", cfg.StatPort)
	log.Printf("   - AI Service:           http://localhost:%s", cfg.AIPort)
	log.Printf("   - Notification Service: http://localhost:%s", cfg.NotificationPort)
	log.Printf("倹 Gateway URL: http://localhost:%s", cfg.GatewayPort)
	log.Printf("当 API Documentation: http://localhost:%s/api/docs", cfg.GatewayPort)

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

	// Statistics Service routes
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

	// Notification Service routes (for future implementation)
	notificationProxy := proxy.NewServiceProxy("http://localhost:" + cfg.NotificationPort)
	notificationRoutes := r.Group("/api/notifications")
	{
		notificationRoutes.Any("/health", notificationProxy.ProxyRequest)
		notificationRoutes.Any("/api/v1/notifications/*path", notificationProxy.ProxyRequest)
	}
}

// setupUtilityRoutes adds utility endpoints
func setupUtilityRoutes(r *gin.Engine, cfg *config.Config) {
	// API documentation endpoint
	r.GET("/api/docs", func(c *gin.Context) {
		docs := map[string]interface{}{
			"title":       "DailyTrackr API Documentation",
			"version":     "1.0.0",
			"description": "Comprehensive productivity tracking API with AI insights",
			"base_url":    "http://localhost:" + cfg.GatewayPort,
			"authentication": map[string]string{
				"type":   "Bearer Token",
				"header": "Authorization: Bearer <jwt_token>",
				"login":  "POST /api/users/auth/login",
			},
			"endpoints": map[string]interface{}{
				"authentication": map[string]string{
					"register": "POST /api/users/auth/register",
					"login":    "POST /api/users/auth/login",
					"profile":  "GET /api/users/api/v1/users/profile",
				},
				"activities": map[string]string{
					"create":       "POST /api/activities/api/v1/activities/",
					"list":         "GET /api/activities/api/v1/activities/",
					"get":          "GET /api/activities/api/v1/activities/:id",
					"update":       "PUT /api/activities/api/v1/activities/:id",
					"delete":       "DELETE /api/activities/api/v1/activities/:id",
					"upload_photo": "POST /api/activities/api/v1/activities/:id/photo",
				},
				"habits": map[string]string{
					"create":     "POST /api/habits/api/v1/habits/",
					"list":       "GET /api/habits/api/v1/habits/",
					"get":        "GET /api/habits/api/v1/habits/:id",
					"update":     "PUT /api/habits/api/v1/habits/:id",
					"delete":     "DELETE /api/habits/api/v1/habits/:id",
					"create_log": "POST /api/habits/api/v1/habits/:id/logs/",
					"get_logs":   "GET /api/habits/api/v1/habits/:id/logs/",
					"get_stats":  "GET /api/habits/api/v1/habits/:id/stats",
				},
				"statistics": map[string]string{
					"dashboard":        "GET /api/stats/api/v1/stats/dashboard",
					"activity_summary": "GET /api/stats/api/v1/stats/activities/summary",
					"activity_chart":   "GET /api/stats/api/v1/stats/activities/chart",
					"habit_progress":   "GET /api/stats/api/v1/stats/habits/progress",
					"expense_report":   "GET /api/stats/api/v1/stats/expenses/report",
				},
				"ai": map[string]string{
					"daily_summary":        "POST /api/ai/api/v1/ai/daily-summary",
					"habit_recommendation": "GET /api/ai/api/v1/ai/habit-recommendation",
					"insights":             "GET /api/ai/api/v1/ai/insights",
					"analyze_activities":   "GET /api/ai/api/v1/ai/analyze-activities",
					"productivity_tips":    "GET /api/ai/api/v1/ai/productivity-tips",
				},
			},
			"examples": map[string]interface{}{
				"login": map[string]interface{}{
					"url":  "POST /api/users/auth/login",
					"body": `{"email": "test123@example.com", "password": "password123"}`,
				},
				"dashboard": map[string]interface{}{
					"url":     "GET /api/stats/api/v1/stats/dashboard",
					"headers": `{"Authorization": "Bearer <jwt_token>"}`,
				},
				"ai_summary": map[string]interface{}{
					"url":     "POST /api/ai/api/v1/ai/daily-summary?date=2025-07-13",
					"headers": `{"Authorization": "Bearer <jwt_token>"}`,
				},
			},
		}

		c.JSON(http.StatusOK, docs)
	})

	// Service status endpoint
	r.GET("/api/status", func(c *gin.Context) {
		services := []string{
			"http://localhost:" + cfg.UserServicePort + "/health",
			"http://localhost:" + cfg.ActivityPort + "/health",
			"http://localhost:" + cfg.HabitPort + "/health",
			"http://localhost:" + cfg.StatPort + "/health",
			"http://localhost:" + cfg.AIPort + "/health",
			"http://localhost:" + cfg.NotificationPort + "/health",
		}

		serviceNames := []string{
			"user-service",
			"activity-service",
			"habit-service",
			"stat-service",
			"ai-service",
			"notification-service",
		}

		status := make(map[string]interface{})
		healthyCount := 0

		for i, serviceURL := range services {
			health := checkServiceHealth(serviceURL)
			status[serviceNames[i]] = map[string]interface{}{
				"status": health,
				"url":    serviceURL,
			}
			if health == "healthy" {
				healthyCount++
			}
		}

		overallStatus := "degraded"
		if healthyCount == len(services) {
			overallStatus = "healthy"
		} else if healthyCount == 0 {
			overallStatus = "down"
		}

		c.JSON(http.StatusOK, gin.H{
			"overall_status":   overallStatus,
			"healthy_services": healthyCount,
			"total_services":   len(services),
			"services":         status,
			"timestamp":        time.Now().Format("2006-01-02 15:04:05"),
		})
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "dailytrackr-gateway",
			"status":  "healthy",
			"version": "1.0.0",
			"uptime":  time.Now().Format("2006-01-02 15:04:05"),
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
