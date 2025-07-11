package proxy

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ServiceProxy handles proxying requests to microservices
type ServiceProxy struct {
	targetURL string
}

// NewServiceProxy creates a new service proxy
func NewServiceProxy(targetURL string) *ServiceProxy {
	return &ServiceProxy{
		targetURL: targetURL,
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
		})
		return
	}

	// Create new request
	reqBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ Failed to read request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// Build target URL with path
	targetPath := c.Request.URL.Path
	// Remove /api/users, /api/activities etc prefix for routing
	if strings.HasPrefix(targetPath, "/api/users") {
		targetPath = strings.TrimPrefix(targetPath, "/api/users")
	} else if strings.HasPrefix(targetPath, "/api/activities") {
		targetPath = strings.TrimPrefix(targetPath, "/api/activities")
	} else if strings.HasPrefix(targetPath, "/api/habits") {
		targetPath = strings.TrimPrefix(targetPath, "/api/habits")
	} else if strings.HasPrefix(targetPath, "/api/ai") {
		targetPath = strings.TrimPrefix(targetPath, "/api/ai")
	}

	if targetPath == "" {
		targetPath = "/health"
	}

	fullURL := target.Scheme + "://" + target.Host + targetPath
	if c.Request.URL.RawQuery != "" {
		fullURL += "?" + c.Request.URL.RawQuery
	}

	log.Printf("🔄 Proxying %s %s -> %s", c.Request.Method, c.Request.URL.Path, fullURL)

	// Create HTTP request
	req, err := http.NewRequest(c.Request.Method, fullURL, bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("❌ Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create proxy request",
		})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request to target service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Failed to proxy request to %s: %v", fullURL, err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"success": false,
			"message": "Service temporarily unavailable",
			"service": sp.targetURL,
		})
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to read service response",
		})
		return
	}

	// Set status and return response
	c.Status(resp.StatusCode)
	c.Writer.Write(respBody)

	log.Printf("✅ Proxy success: %s %s -> %d", c.Request.Method, c.Request.URL.Path, resp.StatusCode)
}
