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

// GetEventsEndpoint retrieves all raw events within a given time range
func GetEventsEndpoint(c *gin.Context) {
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

// PostEventsEndpoint is for testing only
func PostEventsEndpoint(c *gin.Context) {
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
			if err != nil {
				standardAPIResponse(ctx, c, "events.post", err)
				return
			}
		}
	}

	standardAPIResponse(ctx, c, "events.post", err)

}

// PostPredictionsEndpoint is for uploading materialized predictions
func PostPredictionsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	clientID, ok := authenticate(ctx, c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var predictions []types.Prediction
	err := c.BindJSON(&predictions)
	if err == nil {
		for i := range predictions {
			prediction := predictions[i]

			err = backend.StorePrediction(ctx, clientID, &prediction)
			if err != nil {
				standardAPIResponse(ctx, c, "predictions.post", err)
				return
			}
		}
	}

	standardAPIResponse(ctx, c, "predictions.post", err)

}

// SinglePredictionEndpoint returns a single prediction
func SinglePredictionEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	clientID, ok := authenticate(ctx, c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var p types.Prediction
	err := c.BindJSON(&p)
	if err == nil {
		// TODO better auditing
		//logger.Info(ctx, "events.post", "event=%s,type=%s", e.Event, e.EntityType)

		result, err := backend.GetPrediction(ctx, clientID, &p)
		if len(result.Items) == 0 {
			result.Items = make([]types.ItemScore, 0)
		}

		standardJSONResponse(ctx, c, "prediction.single", result, err)
	} else {
		standardAPIResponse(ctx, c, "prediction.single", err)
	}

}
