package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// ModelTrainingEndpoint schedules the training of a model
func ModelTrainingEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "api.model.training"

	// authenticate and authorize
	token := helper.GetToken(ctx, c)
	_, err := authenticateAndAuthorize(ctx, types.ScopeAPIFull, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	// extract values
	modelID := c.Query("id")
	if modelID == "" {
		logger.Warning(ctx, topic, "Empty model ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	err = backend.TrainModel(ctx, modelID)

	if err != nil {
		logger.Warning(ctx, topic, "Issues submitting model for training. Model='%s'. Err=%s", modelID, err.Error())
		helper.StandardAPIResponse(ctx, c, topic, err)
		return
	}

	helper.StandardAPIResponse(ctx, c, topic, err)
}
