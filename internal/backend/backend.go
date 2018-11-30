package backend

import (
	"context"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/ratchetcc/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/pkg/types"
)

// GetEvents queries the events store for events of type event in the time range [start, end]
func GetEvents(ctx context.Context, clientID, event string, start, end int64, page, pageSize int) (*[]EventsStore, error) {
	var events []EventsStore
	var q *datastore.Query

	q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID)

	// filter event type
	if event != "" {
		q = q.Filter("Event =", event)
	}

	// filter time range
	if start > 0 {
		q = q.Filter("Timestamp >=", start)
	}

	if end > 0 {
		q = q.Filter("Timestamp <=", end)
	}

	// order and pageination
	if pageSize > 0 {
		q = q.Order("-Timestamp").Offset((page - 1) * pageSize).Limit(pageSize)
	} else {
		// WARNING: this returns everything !
		q = q.Order("-Timestamp")
	}

	_, err := q.GetAll(ctx, &events)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		events = make([]EventsStore, 0)
	}
	return &events, nil
}

// StoreEvent stores an event in the datastore
func StoreEvent(ctx context.Context, clientID string, event *types.Event) error {

	// deep copy of the struct
	e := EventsStore{
		clientID,
		event.Event,
		event.EntityType,
		event.EntityID,
		event.TargetEntityType,
		event.TargetEntityID,
		event.Properties,
		event.Timestamp,
		util.Timestamp(),
	}

	key := datastore.NewIncompleteKey(ctx, DatastoreEvents, nil)
	_, err := datastore.Put(ctx, key, &e)

	if err != nil {
		logger.Error(ctx, "backend.events.store", err.Error())
	}

	return err
}

// GetPrediction returns a prediction based on a specified model
func GetPrediction(ctx context.Context, clientID string, req *types.Prediction) (*types.Prediction, error) {

	// lookup the model definition
	model, err := GetModel(ctx, clientID, req.Domain)
	if err != nil {
		return nil, err
	}

	// lookup the prediction
	p := types.Prediction{
		EntityID: req.EntityID,
		Domain:   req.Domain,
	}

	key := PredictionKeyString(clientID, model.Domain, req.EntityID, string(model.Revision))
	_, err = memcache.Gob.Get(ctx, key, &p)

	if err != nil {

		ps := PredictionStore{}
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

	ps := PredictionStore{
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

// GetModel returns a model based on the clientID and domain
func GetModel(ctx context.Context, clientID, domain string) (*Model, error) {
	model := Model{}

	// lookup the model definition
	key := "model." + strings.ToLower(clientID) + "." + domain
	_, err := memcache.Gob.Get(ctx, key, &model)

	if err != nil {
		var models []Model
		q := datastore.NewQuery(DatastoreModels).Filter("ClientID =", clientID).Filter("Domain =", domain).Order("-Revision")
		_, err := q.GetAll(ctx, &models)
		if err != nil {
			return nil, err
		}

		if len(models) == 0 {
			return nil, err
		}

		model = models[0]
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = model
			cache.Expiration, _ = time.ParseDuration(ShortCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &model, err
}
