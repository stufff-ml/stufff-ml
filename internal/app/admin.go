package app

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// InitEndpoint creates an initial set of records to get started
func InitEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	token := backend.GetToken(ctx, c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	if token != os.Getenv("ADMIN_TOKEN") {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	err := backend.CreateClientAndAuthentication(ctx, os.Getenv("ADMIN_CLIENT_ID"), os.Getenv("ADMIN_CLIENT_SECRET"), "admin", os.Getenv("ADMIN_CLIENT_TOKEN"))
	if err != nil {
		logger.Error(ctx, "api.seed", err.Error())
	}

	// all good
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
