package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	a "github.com/stufff-ml/stufff-ml/pkg/api"
)

// DefaultEndpoint maps to GET /
func DefaultEndpoint(c *gin.Context) {
	// TODO: real implementation, logging & auditing
	c.JSON(http.StatusOK, gin.H{"vesion": a.Version, "status": "ok"})
}

// RobotsEndpoint maps to GET /robots.txt
func RobotsEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /api/\n"))
}
