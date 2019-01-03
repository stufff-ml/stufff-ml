package app

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/pkg/api"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/jobs"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

//
// TODO better auditing
//

// GetEventsEndpoint retrieves all raw events within a given time range
func GetEventsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "events.get"

	// authenticate and authorize
	token := backend.GetToken(ctx, c)
	clientID, err := backend.AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	// extract values
	event := c.Query("event")
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	if start < 0 {
		start = 0
	}
	end, _ := strconv.Atoi(c.DefaultQuery("end", "0"))
	if end < 0 {
		end = 0
	}
	page, _ := strconv.Atoi(c.DefaultQuery("p", "0"))
	if page < 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("l", "100"))
	if limit < 0 {
		limit = 0
	}

	result, err := backend.GetEvents(ctx, clientID, event, (int64)(start), (int64)(end), page, limit)
	standardJSONResponse(ctx, c, topic, result, err)
}

// PostEventsEndpoint is for testing only
func PostEventsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "events.post"

	// authenticate and authorize
	token := backend.GetToken(ctx, c)
	clientID, err := backend.AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var events []api.Event
	err = c.BindJSON(&events)
	if err == nil {
		// TODO better auditing

		for i := range events {
			e := events[i]
			if e.Timestamp == 0 {
				e.Timestamp = util.Timestamp()
			}
			err = backend.StoreEvent(ctx, clientID, &e)
			if err != nil {
				standardAPIResponse(ctx, c, topic, err)
				return
			}
		}

		logger.Info(ctx, topic, "Received %d event(s).", len(events))
	}

	standardAPIResponse(ctx, c, topic, err)
}

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
				jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, api.JobsBaseURL+"/export?id="+exportID)

				logger.Info(ctx, topic, "Scheduled export of new events. Export='%s'", exportID)
			}
		} else {
			logger.Info(ctx, topic, "Nothing scheduled")
		}
	}

	// logging and standard response
	standardAPIResponse(ctx, c, topic, err)
}

// JobEventsExportEndpoint retrieves all raw events within a given time range
func JobEventsExportEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.events.export"

	// extract values
	exportID := c.Query("id")
	if exportID == "" {
		logger.Warning(ctx, topic, "Empty export ID")
		standardAPIResponse(ctx, c, topic, nil)
		return
	}

	n, err := backend.ExportEvents(ctx, exportID)
	if err != nil {
		logger.Warning(ctx, topic, "Issues exporting new data. Export='%s'. Err=%s", exportID, err.Error())
		standardAPIResponse(ctx, c, topic, err)
		return
	}

	if n > 0 {
		logger.Info(ctx, topic, "Exported new events. Export='%s'", exportID)

		if n == types.ExportBatchSize {
			// more to export, do not merge yet
			jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, api.JobsBaseURL+"/export?id="+exportID)
			logger.Info(ctx, topic, "Re-scheduled export of new events. Export='%s'", exportID)
		} else {
			// schedule merging of files
			jobs.ScheduleJob(ctx, types.BackgroundWorkQueue, api.JobsBaseURL+"/merge?id="+exportID)
			logger.Info(ctx, topic, "Scheduled merge of new events. Export='%s'", exportID)
		}
	}

	standardAPIResponse(ctx, c, topic, err)
}

// JobEventsMergeEndpoint retrieves all exported events files and merges them into one file
func JobEventsMergeEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.events.merge"

	// extract values
	exportID := c.Query("id")
	if exportID == "" {
		logger.Warning(ctx, topic, "Empty export ID")
		standardAPIResponse(ctx, c, topic, nil)
		return
	}

	err := backend.MergeEvents(ctx, exportID)
	if err == nil {
		logger.Info(ctx, topic, "Merged events data. Export='%s'", exportID)
	}

	standardAPIResponse(ctx, c, topic, err)
}
