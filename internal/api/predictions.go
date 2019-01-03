package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	a "github.com/stufff-ml/stufff-ml/pkg/api"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

//
// TODO better auditing
//

// PostPredictionsEndpoint is for uploading materialized predictions
func PostPredictionsEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "predictions.post"

	// authenticate and authorize
	token := GetToken(ctx, c)
	clientID, err := AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var predictions []a.Prediction
	err = c.BindJSON(&predictions)
	if err == nil {
		for i := range predictions {
			prediction := predictions[i]

			err = backend.StorePrediction(ctx, clientID, &prediction)
			if err != nil {
				standardAPIResponse(ctx, c, topic, err)
				return
			}
		}
	}

	standardAPIResponse(ctx, c, topic, err)

}

// GetPredictionEndpoint returns a single prediction
func GetPredictionEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)
	topic := "prediction.single"

	// authenticate and authorize
	token := GetToken(ctx, c)
	clientID, err := AuthenticateAndAuthorize(ctx, "events_access", token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	var p a.Prediction
	err = c.BindJSON(&p)
	if err == nil {
		// TODO better auditing

		result, err := backend.GetPrediction(ctx, clientID, &p)
		if len(result.Items) == 0 {
			result.Items = make([]a.ItemScore, 0)
		}

		standardJSONResponse(ctx, c, topic, result, err)
	} else {
		standardAPIResponse(ctx, c, topic, err)
	}

}
