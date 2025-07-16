package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ServiceProxy handles proxying requests to microservices
type ServiceProxy struct {
	targetURL string
	client    *http.Client
}

// NewServiceProxy creates a new service proxy instance
func NewServiceProxy(targetURL string) *ServiceProxy {
	return &ServiceProxy{
		targetURL: targetURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProxyRequest forwards the request to the target microservice
func (sp *ServiceProxy) ProxyRequest(c *gin.Context) {
	// Build target URL with proper path mapping
	targetPath := sp.buildTargetPathEnhanced(c.Request.URL.Path)
	targetURL := sp.targetURL + targetPath

	// Preserve query parameters
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	log.Printf("🔄 Proxying: %s %s -> %s", c.Request.Method, c.Request.URL.Path, targetURL)

	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
	}

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewReader(bodyBytes))
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create proxy request",
			"error":   err.Error(),
		})
		return
	}

	// Copy headers (excluding connection-specific headers)
	sp.copyHeaders(c.Request.Header, req.Header)

	// Add necessary headers
	req.Header.Set("X-Forwarded-For", c.ClientIP())
	req.Header.Set("X-Forwarded-Proto", "http")
	req.Header.Set("X-Real-IP", c.ClientIP())

	// Execute request
	resp, err := sp.client.Do(req)
	if err != nil {
		log.Printf("❌ Proxy request failed: %v", err)
		c.JSON(http.StatusBadGateway, gin.H{
			"success": false,
			"message": "Service temporarily unavailable",
			"error":   "Failed to connect to backend service",
			"service": sp.targetURL,
		})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	sp.copyResponseHeaders(resp.Header, c.Writer.Header())

	// Copy response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read service response",
			"error":   err.Error(),
		})
		return
	}

	// Set status and return response
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)

	log.Printf("✅ Proxy success: %s %s -> %d (%d bytes)",
		c.Request.Method, c.Request.URL.Path, resp.StatusCode, len(respBody))
}

// buildTargetPathEnhanced builds the target path with enhanced mapping logic
func (sp *ServiceProxy) buildTargetPathEnhanced(originalPath string) string {
	targetPath := originalPath

	// ENHANCED: More precise path mapping based on gateway structure
	switch {
	// User service paths
	case strings.HasPrefix(targetPath, "/api/users/auth/"):
		// /api/users/auth/login -> /auth/login
		targetPath = strings.Replace(targetPath, "/api/users", "", 1)
	case strings.HasPrefix(targetPath, "/api/users/api/v1/users"):
		// /api/users/api/v1/users/profile -> /api/v1/users/profile
		targetPath = strings.Replace(targetPath, "/api/users", "", 1)
	case strings.HasPrefix(targetPath, "/api/users/"):
		// /api/users/health -> /health
		targetPath = strings.Replace(targetPath, "/api/users", "", 1)
	case strings.HasPrefix(targetPath, "/auth/"):
		// Direct auth routes stay as is
		// /auth/login -> /auth/login

	// Activity service paths
	case strings.HasPrefix(targetPath, "/api/activities/api/v1/activities"):
		// /api/activities/api/v1/activities -> /api/v1/activities
		targetPath = strings.Replace(targetPath, "/api/activities", "", 1)
	case strings.HasPrefix(targetPath, "/api/activities/"):
		// /api/activities/health -> /health
		targetPath = strings.Replace(targetPath, "/api/activities", "", 1)

	// Habit service paths
	case strings.HasPrefix(targetPath, "/api/habits/api/v1/habits"):
		// /api/habits/api/v1/habits -> /api/v1/habits
		targetPath = strings.Replace(targetPath, "/api/habits", "", 1)
	case strings.HasPrefix(targetPath, "/api/habits/api/v1/habit-logs"):
		// /api/habits/api/v1/habit-logs -> /api/v1/habit-logs
		targetPath = strings.Replace(targetPath, "/api/habits", "", 1)
	case strings.HasPrefix(targetPath, "/api/habits/"):
		// /api/habits/health -> /health
		targetPath = strings.Replace(targetPath, "/api/habits", "", 1)

	// Statistics service paths
	case strings.HasPrefix(targetPath, "/api/stats/api/v1/stats"):
		// /api/stats/api/v1/stats -> /api/v1/stats
		targetPath = strings.Replace(targetPath, "/api/stats", "", 1)
	case strings.HasPrefix(targetPath, "/api/stats/"):
		// /api/stats/health -> /health
		targetPath = strings.Replace(targetPath, "/api/stats", "", 1)

	// AI service paths
	case strings.HasPrefix(targetPath, "/api/ai/api/v1/ai"):
		// /api/ai/api/v1/ai -> /api/v1/ai
		targetPath = strings.Replace(targetPath, "/api/ai", "", 1)
	case strings.HasPrefix(targetPath, "/api/ai/"):
		// /api/ai/health -> /health
		targetPath = strings.Replace(targetPath, "/api/ai", "", 1)

	// Notification service paths
	case strings.HasPrefix(targetPath, "/api/notifications/api/v1/notifications"):
		// /api/notifications/api/v1/notifications -> /api/v1/notifications
		targetPath = strings.Replace(targetPath, "/api/notifications", "", 1)
	case strings.HasPrefix(targetPath, "/api/notifications/"):
		// /api/notifications/health -> /health
		targetPath = strings.Replace(targetPath, "/api/notifications", "", 1)
	}

	// Default to health check if path is empty
	if targetPath == "" || targetPath == "/" {
		targetPath = "/health"
	}

	return targetPath
}

// copyHeaders copies headers from source to destination with filtering
func (sp *ServiceProxy) copyHeaders(src http.Header, dst http.Header) {
	// Headers to skip (will be set by the target service or client)
	skipHeaders := map[string]bool{
		"Connection":          true,
		"Proxy-Connection":    true,
		"Proxy-Authenticate":  true,
		"Proxy-Authorization": true,
		"Te":                  true,
		"Trailer":             true,
		"Transfer-Encoding":   true,
		"Upgrade":             true,
		"Host":                true, // Let Go set the host
	}

	for key, values := range src {
		normalizedKey := http.CanonicalHeaderKey(key)
		if !skipHeaders[normalizedKey] {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}
}

// copyResponseHeaders copies response headers with filtering
func (sp *ServiceProxy) copyResponseHeaders(src http.Header, dst http.Header) {
	// Headers to skip in response
	skipHeaders := map[string]bool{
		"Connection":        true,
		"Transfer-Encoding": true,
		"Upgrade":           true,
		"Server":            true, // Let Gin set the server header
	}

	for key, values := range src {
		normalizedKey := http.CanonicalHeaderKey(key)
		if !skipHeaders[normalizedKey] {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}

	// Ensure CORS headers are preserved
	dst.Set("Access-Control-Allow-Origin", "*")
	dst.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	dst.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
}
