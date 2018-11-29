package backend

import (
	"context"

	"google.golang.org/appengine/datastore"

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
func RetrieveEvents(ctx context.Context, clientID, event string, start, end int64) (*[]EventsStore, error) {
	var events []EventsStore
	var q *datastore.Query

	if event == "" {
		if start > 0 {
			if end > 0 {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Timestamp >=", start).Filter("Timestamp <=", end).Order("-Timestamp")
			} else {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Timestamp >=", start).Order("-Timestamp")
			}
		} else {
			if end > 0 {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Timestamp <=", end).Order("-Timestamp")
			} else {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Order("-Timestamp")
			}
		}
	} else {
		if start > 0 {
			if end > 0 {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Event =", event).Filter("Timestamp >=", start).Filter("Timestamp <=", end).Order("-Timestamp")
			} else {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Event =", event).Filter("Timestamp >=", start).Order("-Timestamp")
			}
		} else {
			if end > 0 {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Event =", event).Filter("Timestamp <=", end).Order("-Timestamp")
			} else {
				q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Event =", event).Order("-Timestamp")
			}
		}
	}

	_, err := q.GetAll(ctx, &events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}
