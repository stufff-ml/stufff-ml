package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/stufff-ml/stufff-ml/internal/backend"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/gae/logger"
)

// standardAPIResponse is the default way to respond to API requests
func standardAPIResponse(ctx context.Context, c *gin.Context, topic string, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		logger.Error(ctx, "api.response", err.Error())
		// TODO proper error handling. For now 400 it is
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}

// standardJSONResponse is the default way to respond to API requests
func standardJSONResponse(ctx context.Context, c *gin.Context, topic string, res interface{}, err error) {
	if err == nil {
		if res == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusOK, res)
		}
	} else {
		logger.Error(ctx, "api.response", err.Error())
		// TODO proper error handling. For now 400 it is
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}

// authorizeRequest looks for a token and returns the matching model ID
func authenticate(ctx context.Context, c *gin.Context) (string, bool) {

	auth := c.Request.Header["Authorization"]
	if len(auth) == 0 {
		return "", false
	}

	parts := strings.Split(auth[0], " ")
	if len(parts) != 2 {
		return "", false
	}

	return backend.ClientIDFromToken(ctx, parts[1])
}
