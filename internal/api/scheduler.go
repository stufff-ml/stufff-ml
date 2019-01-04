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

	"github.com/stufff-ml/stufff-ml/internal/jobs"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// ScheduleEventsExportEndpoint retrieves all raw events within a given time range
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
				jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, a.JobsBaseURL+"/export?id="+exportID)

				logger.Info(ctx, topic, "Scheduled export of new events. Export='%s'", exportID)
			}
		} else {
			logger.Info(ctx, topic, "Nothing scheduled")
		}
	}

	// logging and standard response
	helper.StandardAPIResponse(ctx, c, topic, err)
}
