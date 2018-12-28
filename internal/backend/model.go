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
func CreateModel(ctx context.Context, clientID, name string) (*ModelDS, error) {

	model := ModelDS{
		ClientID:         clientID,
		Name:             name,
		Revision:         1,
		TrainingSchedule: 60,
		NextSchedule:     0,
		Created:          util.Timestamp(),
	}

	key := ModelKey(ctx, clientID, name)
	_, err := datastore.Put(ctx, key, &model)
	if err != nil {
		logger.Error(ctx, "backend.model.create", err.Error())
		return nil, err
	}

	return &model, nil

}

// GetModel returns a model based on the clientID and domain
func GetModel(ctx context.Context, clientID, name string) (*ModelDS, error) {
	model := ModelDS{}

	// lookup the model definition
	key := "model." + strings.ToLower(clientID+"."+name)
	_, err := memcache.Gob.Get(ctx, key, &model)

	if err != nil {
		var models []ModelDS
		q := datastore.NewQuery(DatastoreModels).Filter("ClientID =", clientID).Filter("Name =", name).Order("-Revision")
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
