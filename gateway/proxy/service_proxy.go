package proxy

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ServiceProxy handles proxying requests to microservices
type ServiceProxy struct {
	targetURL string
	client    *http.Client
}

// NewServiceProxy creates a new service proxy with optimized settings
func NewServiceProxy(targetURL string) *ServiceProxy {
	return &ServiceProxy{
		targetURL: targetURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     30 * time.Second,
			},
		},
	}
}

// ProxyRequest forwards the request to the target service
func (sp *ServiceProxy) ProxyRequest(c *gin.Context) {
	// Parse target URL
	target, err := url.Parse(sp.targetURL)
	if err != nil {
		log.Printf("❌ Failed to parse target URL %s: %v", sp.targetURL, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gateway configuration error",
			"error":   "Invalid target URL",
		})
		return
	}

	// Read request body
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	// Build target URL with path
	targetPath := sp.buildTargetPath(c.Request.URL.Path)
	fullURL := target.Scheme + "://" + target.Host + targetPath

	// Add query parameters
	if c.Request.URL.RawQuery != "" {
		fullURL += "?" + c.Request.URL.RawQuery
	}

	log.Printf("🔄 Proxying %s %s -> %s", c.Request.Method, c.Request.URL.Path, fullURL)

	// Create HTTP request with context for timeout
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, c.Request.Method, fullURL, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create proxy request",
			"error":   err.Error(),
		})
		return
	}

	// Copy headers with filtering
	sp.copyHeaders(c.Request.Header, req.Header)

	// Make request to target service
	resp, err := sp.client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to proxy request to %s: %v", fullURL, err)

		// Enhanced error response based on error type
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{
				"success": false,
				"message": "Service timeout",
				"service": sp.targetURL,
				"error":   "Request timed out",
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"success": false,
				"message": "Service temporarily unavailable",
				"service": sp.targetURL,
				"error":   err.Error(),
			})
		}
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

// buildTargetPath builds the target path by removing gateway prefixes
func (sp *ServiceProxy) buildTargetPath(originalPath string) string {
	targetPath := originalPath

	// FIXED: Better path mapping logic
	switch {
	case strings.HasPrefix(targetPath, "/api/users/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/users")
	case strings.HasPrefix(targetPath, "/api/activities/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/activities")
	case strings.HasPrefix(targetPath, "/api/habits/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/habits")
	case strings.HasPrefix(targetPath, "/api/stats/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/stats")
	case strings.HasPrefix(targetPath, "/api/ai/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/ai")
	case strings.HasPrefix(targetPath, "/api/notifications/"):
		targetPath = strings.TrimPrefix(targetPath, "/api/notifications")
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
	}

	for key, values := range src {
		if !skipHeaders[key] {
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
	}

	for key, values := range src {
		if !skipHeaders[key] {
			for _, value := range values {
				dst.Add(key, value)
			}
		}
	}
}
