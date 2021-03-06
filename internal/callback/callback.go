package callback

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
	"github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"
)

// ModelTrainingEndpoint is used to receive notifications on completed model training
func ModelTrainingEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "callback.model.training"

	// extract values
	clientID := c.Query("id")
	if clientID == "" {
		logger.Warning(ctx, topic, "Empty client ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	jobID := c.Query("job")
	if jobID == "" {
		logger.Warning(ctx, topic, "Missing job ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	status := c.Query("status")
	if status == "" {
		logger.Warning(ctx, topic, "Missing status")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	if status == "ok" {
		// schedule import
		q := fmt.Sprintf("%s/import?id=%s", api.JobsPrefix, jobID)
		backend.ScheduleJob(ctx, types.BackgroundWorkQueue, q)

		logger.Info(ctx, topic, "Scheduled import of training data. Query='%s'", q)
	}

	err := backend.MarkModelTrainingDone(ctx, jobID, status)

	// logging and standard response
	logger.Info(ctx, topic, "Training job '%s' finished. Client ID=%s", jobID, clientID)
	helper.StandardAPIResponse(ctx, c, topic, err)
}
