package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// DebugEndpoint is for testing only
func DebugEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// MigrateEndpoint is for testing only
func MigrateEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// SeedEndpoint creates an initial set of records to get started
func SeedEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	err := backend.CreateClientAndAuthentication(ctx, "aaaa", "aaaa", "xoxo-ffffffff")
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	_, err = backend.CreateModel(ctx, "aaaa", "buy", "buy")
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	// all good
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
