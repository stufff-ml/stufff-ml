package backend

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// ModelKey key on collection MODELS
func ModelKey(ctx context.Context, modelID string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreModels, modelID, 0, nil)
}

// ClientResourceKey key on collection CLIENT_RESOURCES
func ClientResourceKey(ctx context.Context, clientID string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreClientResources, clientID, 0, nil)
}

// AuthorizationKey key on collection AUTHORIZATIONS
func AuthorizationKey(ctx context.Context, token string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreAuthorizations, token, 0, nil)
}
