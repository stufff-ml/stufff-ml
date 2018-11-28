package backend

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// AuthorizationKey creates a datastore key for a workspace authorization based on the team_id.
func AuthorizationKey(ctx context.Context, id string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreEvents, id, 0, nil)
}
