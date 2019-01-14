package backend

import (
	"context"

	"github.com/stufff-ml/stufff-ml/pkg/api"
)

// GetPrediction returns a prediction based on a specified model
func GetPrediction(ctx context.Context, clientID string, req *api.Prediction) (*api.Prediction, error) {

	// lookup the prediction
	p := api.Prediction{
		EntityID: req.EntityID,
		Domain:   req.Domain,
		Items:    make([]api.ItemScore, 0),
	}

	return &p, nil
}

// StorePrediction stores a materialized prediction in the datastore
func StorePrediction(ctx context.Context, clientID string, prediction *api.Prediction) error {

	return nil
}
