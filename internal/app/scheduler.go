package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/jobs"
	"github.com/stufff-ml/stufff-ml/pkg/types"
)

//
// TODO better auditing
//

// ScheduleEventsExportEndpoint retrieves all raw events within a given time range
func ScheduleEventsExportEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	var models []backend.Model
	now := util.Timestamp()

	q := datastore.NewQuery(backend.DatastoreModels).Filter("NextSchedule <=", now)

	_, err := q.GetAll(ctx, &models)
	if err == nil {
		if len(models) > 0 {
			for i := range models {
				modelID := fmt.Sprintf("%s.%s", models[i].ClientID, models[i].Domain)
				jobs.ScheduleJob(ctx, backend.BackgroundWorkQueue, types.JobsBaseURL+"/export?id="+modelID)

				logger.Info(ctx, "scheduler.events.export", "Scheduled export of model '%s'.", modelID)
			}
		} else {
			logger.Info(ctx, "scheduler.events.export", "Nothing scheduled.")
		}
	}

	// logging and standard response
	standardAPIResponse(ctx, c, "scheduler.events.export", err)
}
