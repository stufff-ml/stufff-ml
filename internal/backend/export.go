package backend

import (
	"context"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/types"
)

// CreateDefaultExport creates an initial export definition
func CreateDefaultExport(ctx context.Context, clientID string) (*types.ExportDS, error) {
	return CreateExport(ctx, clientID, "default")
}

// CreateExport creates a new export definition
func CreateExport(ctx context.Context, clientID, event string) (*types.ExportDS, error) {

	model := types.ExportDS{
		ClientID:       clientID,
		Event:          event,
		Exported:       0,
		ExportedTotal:  0,
		ExportSchedule: 15,
		NextSchedule:   0,
		Created:        util.Timestamp(),
	}

	key := ExportKey(ctx, clientID, event)
	_, err := datastore.Put(ctx, key, &model)
	if err != nil {
		logger.Error(ctx, "backend.export.create", err.Error())
		return nil, err
	}

	return &model, nil
}

// GetExport returns an export definition based on the clientID and event
func GetExport(ctx context.Context, clientID, event string) (*types.ExportDS, error) {
	var export types.ExportDS

	// lookup the model definition
	key := "export." + strings.ToLower(clientID+"."+event)
	_, err := memcache.Gob.Get(ctx, key, &export)

	if err != nil {
		var exports []types.ExportDS
		q := datastore.NewQuery(types.DatastoreExports).Filter("ClientID =", clientID).Filter("Event =", event)
		_, err := q.GetAll(ctx, &exports)
		if err != nil {
			return nil, err
		}

		if len(exports) == 0 {
			return nil, err
		}

		export = exports[0]
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = export
			cache.Expiration, _ = time.ParseDuration(types.ShortCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &export, nil
}

// markExported writes an export record back to the datastore with updated metadata
func markExported(ctx context.Context, clientID, event string, exported int, ts, next int64) error {
	var export types.ExportDS

	key := ExportKey(ctx, clientID, event)
	err := datastore.Get(ctx, key, &export)
	if err != nil {
		return err
	}

	export.Exported = exported
	export.ExportedTotal += exported

	export.LastExported = ts
	export.NextSchedule = next

	_, err = datastore.Put(ctx, key, &export)
	if err != nil {
		return err
	}

	// invalidate the cache
	ckey := "export." + strings.ToLower(clientID+"."+event)
	err = memcache.Delete(ctx, ckey)

	return err
}
