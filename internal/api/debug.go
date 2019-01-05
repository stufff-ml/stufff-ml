package api

import (
	"net/http"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/pkg/helper"
	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"
)

// DebugEndpoint is for testing only
func DebugEndpoint(c *gin.Context) {
	topic := "DEBUG"
	ctx := appengine.NewContext(c.Request)

	err := backend.TrainModel(ctx, "869746c7855d.default")

	helper.StandardAPIResponse(ctx, c, topic, err)
}

// MigrateEndpoint is for testing only
func MigrateEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
