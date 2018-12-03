package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DebugEndpoint is for testing only
func DebugEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MigrateEndpoint is for testing only
func MigrateEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
