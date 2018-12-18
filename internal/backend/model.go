package backend

import (
	"context"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
)

// CreateModel creates an initial model definition
func CreateModel(ctx context.Context, clientID, domain string) (*Model, error) {

	model := Model{
		ClientID:         clientID,
		Domain:           domain,
		Revision:         1,
		ExportSchedule:   15,
		TrainingSchedule: 60,
		NextSchedule:     0,
		LastExported:     0,
		Created:          util.Timestamp(),
	}

	key := ModelKey(ctx, clientID, domain)
	_, err := datastore.Put(ctx, key, &model)
	if err != nil {
		logger.Error(ctx, "backend.model.create", err.Error())
		return nil, err
	}

	return &model, nil

}

// GetModel returns a model based on the clientID and domain
func GetModel(ctx context.Context, clientID, domain string) (*Model, error) {
	model := Model{}

	// lookup the model definition
	key := "model." + strings.ToLower(clientID+"."+domain)
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

	return &model, nil
}

// MarkModelExported writes a model record back to the datastore with updated metadata
func MarkModelExported(ctx context.Context, clientID, domain string, exported, next int64) error {
	var model Model

	key := ModelKey(ctx, clientID, domain)
	err := datastore.Get(ctx, key, &model)
	if err != nil {
		return err
	}

	model.LastExported = exported
	model.NextSchedule = next

	_, err = datastore.Put(ctx, key, &model)
	if err != nil {
		return err
	}

	// invalidate the cache
	ckey := "model." + strings.ToLower(clientID+"."+domain)
	err = memcache.Delete(ctx, ckey)

	return err
}
