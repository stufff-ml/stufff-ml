package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/gae/logger"
	"google.golang.org/appengine"

	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// GetPredictionEndpoint returns a single prediction
func GetPredictionEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "api.predictions.get"

	// authenticate and authorize
	token := helper.GetToken(ctx, c)
	clientID, err := authenticateAndAuthorize(ctx, types.ScopeUserFull, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	// extract values
	entityType := c.Query("type")
	if entityType == "" {
		logger.Warning(ctx, topic, "Empty entity type")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	targetEntityType := c.Query("target")
	if targetEntityType == "" {
		logger.Warning(ctx, topic, "Empty target entity type")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	entityID := c.Query("id")
	if entityID == "" {
		logger.Warning(ctx, topic, "Empty entity ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	model := c.Query("model")
	if model == "" {
		model = "default"
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("l", "10"))
	if limit < 0 {
		limit = 0
	}

	result, err := backend.GetPrediction(ctx, clientID, model, entityID, entityType, targetEntityType, limit)
	helper.StandardJSONResponse(ctx, c, topic, result, err)

}
