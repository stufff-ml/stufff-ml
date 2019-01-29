package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	a "github.com/stufff-ml/stufff-ml/pkg/api"
	"github.com/stufff-ml/stufff-ml/pkg/helper"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/internal/types"
)

// GetEventsEndpoint retrieves all raw events within a given time range
func GetEventsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "api.events.get"

	// authenticate and authorize
	token := helper.GetToken(ctx, c)
	clientID, err := authenticateAndAuthorize(ctx, types.ScopeUserFull, token)
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
	helper.StandardJSONResponse(ctx, c, topic, result, err)
}

// PostEventsEndpoint is for testing only
func PostEventsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "api.events.post"

	// authenticate and authorize
	token := helper.GetToken(ctx, c)
	clientID, err := authenticateAndAuthorize(ctx, types.ScopeUserFull, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var events []a.Event
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
				helper.StandardAPIResponse(ctx, c, topic, err)
				return
			}
		}

		logger.Info(ctx, topic, "Received %d event(s).", len(events))
	}

	helper.StandardAPIResponse(ctx, c, topic, err)
}
