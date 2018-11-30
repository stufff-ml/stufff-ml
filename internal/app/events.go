package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/pkg/types"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// GetEventEndpoint retrieves all raw events within a given time range
func GetEventEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	clientID, ok := authenticate(ctx, c)
	if !ok {
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

	result, err := backend.RetrieveEvents(ctx, clientID, event, (int64)(start), (int64)(end), page, pageSize)
	standardJSONResponse(ctx, c, "events.get", result, err)

}

// PostEventEndpoint is for testing only
func PostEventEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	clientID, ok := authenticate(ctx, c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var events []types.Event
	err := c.BindJSON(&events)
	if err == nil {
		// TODO better auditing
		//logger.Info(ctx, "events.post", "event=%s,type=%s", e.Event, e.EntityType)

		for i := range events {
			e := events[i]
			if e.Timestamp == 0 {
				e.Timestamp = util.Timestamp()
			}
			err = backend.StoreEvent(ctx, clientID, &e)
		}
	}

	standardAPIResponse(ctx, c, "events.post", err)

}
