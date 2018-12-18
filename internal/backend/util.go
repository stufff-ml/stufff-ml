package backend

import (
	"fmt"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
)

// ModelKey key on collection MODELS
func ModelKey(ctx context.Context, clientID, domain string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreModels, clientID+"."+domain, 0, nil)
}

// ClientResourceKey key on collection CLIENT_RESOURCES
func ClientResourceKey(ctx context.Context, clientID string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreClientResources, clientID, 0, nil)
}

// AuthorizationKey key on collection AUTHORIZATIONS
func AuthorizationKey(ctx context.Context, token string) *datastore.Key {
	return datastore.NewKey(ctx, DatastoreAuthorizations, token, 0, nil)
}

// PredictionKeyString returns the composite key string for a prediction
func PredictionKeyString(clientID, domain, entityID, revision string) string {
	return clientID + "." + domain + "." + entityID + "." + revision
}

// PredictionKey key on collection PREDICTIONS
func PredictionKey(ctx context.Context, k string) *datastore.Key {
	return datastore.NewKey(ctx, DatastorePredictions, k, 0, nil)
}

func (e *EventsStore) ToCSV() string {
	if len(e.Properties) == 0 {
		return fmt.Sprintf("%s,%s,%s,%s,%s,%d,''\n", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%d,'%s'\n", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp, strings.Join(e.Properties, ","))
}
