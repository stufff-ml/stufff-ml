package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/stufff-ml/stufff-ml/internal/backend"
	"github.com/stufff-ml/stufff-ml/pkg/types"
)

// SinglePredictionEndpoint is for testing only
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

		result, err := backend.Prediction(ctx, clientID, &p)
		if len(result.Items) == 0 {
			result.Items = make([]types.ItemScore, 0)
		}

		standardJSONResponse(ctx, c, "prediction.single", result, err)
	} else {
		standardAPIResponse(ctx, c, "prediction.single", err)
	}

}
