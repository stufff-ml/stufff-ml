package api

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	a "github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// ScheduleEventsExportEndpoint schedules the export of new events
func ScheduleEventsExportEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "scheduler.events.export"

	var exports []types.ExportDS
	now := util.Timestamp()

	q := datastore.NewQuery(types.DatastoreExports).Filter("NextSchedule <=", now)
	_, err := q.GetAll(ctx, &exports)

	if err == nil {
		if len(exports) > 0 {
			for i := range exports {
				exportID := fmt.Sprintf("%s.%s", exports[i].ClientID, exports[i].Event)
				backend.ScheduleJob(ctx, types.BackgroundWorkQueue, a.JobsBaseURL+"/export?id="+exportID)

				logger.Info(ctx, topic, "Scheduled export of new events. Export='%s'", exportID)
			}
		} else {
			logger.Info(ctx, topic, "Nothing scheduled")
		}
	}

	// logging and standard response
	helper.StandardAPIResponse(ctx, c, topic, err)
}

// ScheduleModelTrainingEndpoint schedules periodic model training
func ScheduleModelTrainingEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "scheduler.model.training"

	var models []types.ModelDS
	now := util.Timestamp()

	q := datastore.NewQuery(types.DatastoreModels).Filter("NextSchedule <=", now)
	_, err := q.GetAll(ctx, &models)

	if err == nil {
		if len(models) > 0 {
			for i := range models {
				modelID := fmt.Sprintf("%s.%s", models[i].ClientID, models[i].Name)
				backend.ScheduleJob(ctx, types.BackgroundWorkQueue, a.JobsBaseURL+"/train?id="+modelID)

				logger.Info(ctx, topic, "Scheduled model training. Model='%s'", modelID)
			}
		} else {
			logger.Info(ctx, topic, "Nothing scheduled")
		}
	}

	// logging and standard response
	helper.StandardAPIResponse(ctx, c, topic, err)
}
