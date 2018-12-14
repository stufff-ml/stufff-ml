package app

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"cloud.google.com/go/storage"

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
	standardJSONResponse(ctx, c, "events.get", result, err)

}

// PostEventsEndpoint is for testing only
func PostEventsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

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
		//logger.Info(ctx, "events.post", "event=%s,type=%s", e.Event, e.EntityType)

		for i := range events {
			e := events[i]
			if e.Timestamp == 0 {
				e.Timestamp = util.Timestamp()
			}
			err = backend.StoreEvent(ctx, clientID, &e)
			if err != nil {
				standardAPIResponse(ctx, c, "events.post", err)
				return
			}
		}
	}

	standardAPIResponse(ctx, c, "events.post", err)

}

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

// JobEventsExportEndpoint retrieves all raw events within a given time range
func JobEventsExportEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// extract values
	modelID := c.Query("id")
	if modelID == "" {
		logger.Warning(ctx, "jobs.events.export", "Empty model ID")
		standardAPIResponse(ctx, c, "jobs.events.export", nil)
		return
	}

	p := strings.Split(modelID, ".")
	clientID := p[0]
	domain := p[1]

	model, err := backend.GetModel(ctx, clientID, domain)
	if err != nil {
		logger.Warning(ctx, "jobs.events.export", "Model not found: %s", modelID)
		standardAPIResponse(ctx, c, "jobs.events.export", err)
		return
	}

	// timerange: ]start, end]
	start := model.LastExported
	end := util.Timestamp()

	// export stuff
	var events *[]backend.EventsStore
	events, err = backend.GetEvents(ctx, model.ClientID, "", start, end, 0, 0)
	if err != nil {
		logger.Warning(ctx, "jobs.events.export", "Could not query events for model: %s", modelID)
		standardAPIResponse(ctx, c, "jobs.events.export", err)
		return
	}

	// only if there is something to export
	if len(*events) > 0 {

		// create a file on Cloud Storage
		client, err := storage.NewClient(ctx)
		if err != nil {
			logger.Warning(ctx, "jobs.events.export", "Could not write to storage %s", modelID)
			standardAPIResponse(ctx, c, "jobs.events.export", err)
			return
		}

		bucket := client.Bucket("exports.stufff.review")

		fileName := fmt.Sprintf("%s/%s_%d.csv", modelID, modelID, end)
		w := bucket.Object(fileName).NewWriter(ctx)
		w.ContentType = "text/plain"
		defer w.Close()

		// write to file
		for i := range *events {
			w.Write([]byte(backend.EventStoreToString(&(*events)[i]) + "\n"))
		}

		logger.Info(ctx, "jobs.events.export", "Wrote %d events to file '%s'", len(*events), fileName)
	}

	// uodate metadata
	model.LastExported = end
	model.NextSchedule = util.IncT(end, model.TrainingSchedule)
	err = backend.MarkModelExported(ctx, clientID, domain, end, util.IncT(end, model.TrainingSchedule))
	if err != nil {
		logger.Warning(ctx, "jobs.events.export", "Could not update metadata for model: %s", modelID)
		standardAPIResponse(ctx, c, "jobs.events.export", err)
		return
	}

	// logging and standard response
	logger.Info(ctx, "jobs.events.export", "Exported data for models '%s'", modelID)
	standardAPIResponse(ctx, c, "jobs.events.export", nil)
}
