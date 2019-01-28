package api

import (
	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/pkg/helper"
)

// ModelTrainingCallback is used to receive notifications on completed model training
func ModelTrainingCallback(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "callback.model.training"

	// extract values
	clientID := c.Query("id")
	if clientID == "" {
		logger.Warning(ctx, topic, "Empty client ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	job := c.Query("job")
	if job == "" {
		logger.Warning(ctx, topic, "Missing job ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	// logging and standard response
	logger.Info(ctx, topic, "Training job '%s' finished. Client ID=%s", job, clientID)
	helper.StandardAPIResponse(ctx, c, topic, nil)
}
