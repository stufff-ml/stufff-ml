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

	"github.com/stufff-ml/stufff-ml/pkg/types"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/jobs"
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
	page, _ := strconv.Atoi(c.DefaultQuery("p", "1"))
	if page < 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(c.DefaultQuery("s", "0"))
	if pageSize < 0 {
		pageSize = 0
	}

	result, err := backend.GetEvents(ctx, clientID, event, (int64)(start), (int64)(end), page, pageSize)
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

	var events []types.Event
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

	var models []backend.Model
	now := util.Timestamp()

	q := datastore.NewQuery(backend.DatastoreModels).Filter("NextSchedule <=", now)

	_, err := q.GetAll(ctx, &models)
	if err == nil {
		if len(models) > 0 {
			for i := range models {
				modelID := fmt.Sprintf("%s.%s", models[i].ClientID, models[i].Domain)
				jobs.ScheduleJob(ctx, backend.BackgroundWorkQueue, types.JobsBaseURL+"/export?id="+modelID)

				logger.Info(ctx, topic, "Scheduled export of new events. Model='%s'", modelID)
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
	modelID := c.Query("id")
	if modelID == "" {
		logger.Warning(ctx, topic, "Empty model ID")
		standardAPIResponse(ctx, c, topic, nil)
		return
	}

	n, err := backend.ExportEvents(ctx, modelID)
	if err != nil {
		logger.Warning(ctx, topic, "Issues exporting new data. Model='%s'. Err=%s", modelID, err.Error())
		standardAPIResponse(ctx, c, topic, err)
		return
	}

	if n > 0 {
		logger.Info(ctx, topic, "Exported new data. Model='%s'", modelID)

		// schedule merging of files
		jobs.ScheduleJob(ctx, backend.BackgroundWorkQueue, types.JobsBaseURL+"/merge?id="+modelID)
		logger.Info(ctx, topic, "Scheduled merge of new events. Model='%s'", modelID)
	}

	standardAPIResponse(ctx, c, topic, err)
}

// JobEventsMergeEndpoint retrieves all exported events files and merges them into one file
func JobEventsMergeEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "jobs.events.merge"

	// extract values
	modelID := c.Query("id")
	if modelID == "" {
		logger.Warning(ctx, topic, "Empty model ID")
		standardAPIResponse(ctx, c, topic, nil)
		return
	}

	err := backend.MergeEvents(ctx, modelID)
	if err == nil {
		logger.Info(ctx, topic, "Merged events data. Model='%s'", modelID)
	}

	standardAPIResponse(ctx, c, topic, err)
}
