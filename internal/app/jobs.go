package app

import (
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// ExportEventsJobEndpoint retrieves all raw events within a given time range
func ExportEventsJobEndpoint(c *gin.Context) {
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

	for i := range *events {
		logger.Info(ctx, "jobs.events.export", "event=%s, type=%s", (*events)[i].Event, (*events)[i].EntityType)
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
