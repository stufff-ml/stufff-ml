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

// RetrieveEvents queries the events store for events of type event in the time range [start, end]
func RetrieveEvents(ctx context.Context, clientID, event string, start, end int64, page, pageSize int) (*[]EventsStore, error) {
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

// Prediction returns a prediction based on a specified model
func Prediction(ctx context.Context, clientID string, req *types.Prediction) (*types.Prediction, error) {
	model := Model{}

	// lookup the model definition
	key := "model." + strings.ToLower(clientID) + "." + req.Domain
	_, err := memcache.Gob.Get(ctx, key, &model)

	if err != nil {
		var models []Model
		q := datastore.NewQuery(DatastoreModels).Filter("ClientID =", clientID).Filter("Domain =", req.Domain).Order("-Revision")
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

	// lookup the prediction
	p := types.Prediction{
		EntityID:   req.EntityID,
		EntityType: req.EntityType,
		Domain:     req.Domain,
	}

	key = PredictionKeyString(clientID, model.ModelID, string(model.Revision), model.Domain, req.EntityID)
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
