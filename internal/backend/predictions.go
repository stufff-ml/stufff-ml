package backend

import (
	"context"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/pkg/types"
)

// GetPrediction returns a prediction based on a specified model
func GetPrediction(ctx context.Context, clientID string, req *types.Prediction) (*types.Prediction, error) {

	// lookup the prediction
	p := types.Prediction{
		EntityID: req.EntityID,
		Domain:   req.Domain,
		Items:    make([]types.ItemScore, 0),
	}

	// lookup the model definition
	model, err := GetModel(ctx, clientID, req.Domain)
	if err != nil {
		return &p, err
	}

	key := PredictionKeyString(clientID, "model.Domain", req.EntityID, string(model.Revision))
	_, err = memcache.Gob.Get(ctx, key, &p)

	if err != nil {

		ps := PredictionDS{}
		err = datastore.Get(ctx, PredictionKey(ctx, key), &ps)
		if err == nil {

			p.Items = ps.Items

			cache := memcache.Item{}
			cache.Key = key
			cache.Object = &p
			cache.Expiration, _ = time.ParseDuration(ShortCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		}
	}

	return &p, nil
}

// StorePrediction stores a materialized prediction in the datastore
func StorePrediction(ctx context.Context, clientID string, prediction *types.Prediction) error {

	model, err := GetModel(ctx, clientID, prediction.Domain)

	ps := PredictionDS{
		ClientID: clientID,
		Domain:   prediction.Domain,
		EntityID: prediction.EntityID,
		Revision: model.Revision,
		Items:    prediction.Items,
		Created:  util.Timestamp(),
	}

	key := PredictionKey(ctx, PredictionKeyString(clientID, prediction.Domain, prediction.EntityID, string(model.Revision)))
	_, err = datastore.Put(ctx, key, &ps)
	if err != nil {
		logger.Error(ctx, "backend.prediction.store", err.Error())
	}

	return err
}
