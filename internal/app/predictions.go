package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/stufff-ml/stufff-ml/pkg/types"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

//
// TODO better auditing
//

// PostPredictionsEndpoint is for uploading materialized predictions
func PostPredictionsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// authenticate and authorize
	token := backend.GetToken(ctx, c)
	clientID, err := backend.AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var predictions []types.Prediction
	err = c.BindJSON(&predictions)
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

// GetPredictionEndpoint returns a single prediction
func GetPredictionEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// authenticate and authorize
	token := backend.GetToken(ctx, c)
	clientID, err := backend.AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var p types.Prediction
	err = c.BindJSON(&p)
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
