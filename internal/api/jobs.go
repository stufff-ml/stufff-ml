package api

import (
	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"

	a "github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/jobs"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// JobEventsExportEndpoint retrieves all raw events within a given time range
func JobEventsExportEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.events.export"

	// extract values
	exportID := c.Query("id")
	if exportID == "" {
		logger.Warning(ctx, topic, "Empty export ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	n, err := backend.ExportEvents(ctx, exportID)
	if err != nil {
		logger.Warning(ctx, topic, "Issues exporting new data. Export='%s'. Err=%s", exportID, err.Error())
		helper.StandardAPIResponse(ctx, c, topic, err)
		return
	}

	if n > 0 {
		logger.Info(ctx, topic, "Exported new events. Export='%s'", exportID)

		if n == types.ExportBatchSize {
			// more to export, do not merge yet
			jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, a.JobsBaseURL+"/export?id="+exportID)
			logger.Info(ctx, topic, "Re-scheduled export of new events. Export='%s'", exportID)
		} else {
			// schedule merging of files
			jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, a.JobsBaseURL+"/merge?id="+exportID)
			logger.Info(ctx, topic, "Scheduled merge of new events. Export='%s'", exportID)
		}
	}

	helper.StandardAPIResponse(ctx, c, topic, err)
}

// JobEventsMergeEndpoint retrieves all exported events files and merges them into one file
func JobEventsMergeEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.events.merge"

	// extract values
	exportID := c.Query("id")
	if exportID == "" {
		logger.Warning(ctx, topic, "Empty export ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	err := backend.MergeEvents(ctx, exportID)
	if err == nil {
		logger.Info(ctx, topic, "Merged events data. Export='%s'", exportID)
	}

	helper.StandardAPIResponse(ctx, c, topic, err)
}

// JobModelTrainingEndpoint schedules the training of a model
func JobModelTrainingEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.model.training"

	// extract values
	modelID := c.Query("id")
	if modelID == "" {
		logger.Warning(ctx, topic, "Empty model ID")
		helper.StandardAPIResponse(ctx, c, topic, nil)
		return
	}

	err := backend.SubmitModel(ctx, modelID)
	if err != nil {
		logger.Warning(ctx, topic, "Issues submitting model for training. Model='%s'. Err=%s", modelID, err.Error())
		helper.StandardAPIResponse(ctx, c, topic, err)
		return
	}

	helper.StandardAPIResponse(ctx, c, topic, err)
}
