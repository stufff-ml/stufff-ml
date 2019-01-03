package backend

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/stufff-ml/stufff-ml/internal/types"
)

// ModelKey key on collection MODELS
func ModelKey(ctx context.Context, clientID, name string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreModels, clientID+"."+name, 0, nil)
}

// ExportKey key on collection EXPORTS
func ExportKey(ctx context.Context, clientID, event string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreExports, clientID+"."+event, 0, nil)
}

// ClientResourceKey key on collection CLIENT_RESOURCES
func ClientResourceKey(ctx context.Context, clientID string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreClientResources, clientID, 0, nil)
}

// AuthorizationKey key on collection AUTHORIZATIONS
func AuthorizationKey(ctx context.Context, token string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreAuthorizations, token, 0, nil)
}

// PredictionKeyString returns the composite key string for a prediction
func PredictionKeyString(clientID, domain, entityID, revision string) string {
	return clientID + "." + domain + "." + entityID + "." + revision
}

// PredictionKey key on collection PREDICTIONS
func PredictionKey(ctx context.Context, k string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastorePredictions, k, 0, nil)
}
