package backend

import (
	"context"
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/stufff-ml/stufff-ml/internal/types"
	"github.com/stufff-ml/stufff-ml/pkg/api"
)

// GetEvents queries the events store for events of type event in the time range [start, end]
func GetEvents(ctx context.Context, clientID, event string, start, end int64, page, limit int) (*[]types.EventDS, error) {
	topic := "backend.events.get"

	var events []types.EventDS
	var q *datastore.Query

	q = datastore.NewQuery(types.DatastoreEvents).Filter("ClientID =", clientID)

	// filter event type
	if event != "" {
		q = q.Filter("Event =", event)
	}

	// filter time range
	if start > 0 {
		q = q.Filter("Timestamp >", start)
	}

	if end > 0 {
		q = q.Filter("Timestamp <=", end)
	}

	// order and pageination
	if page == 0 && limit > 0 {
		q = q.Order("Timestamp").Limit(limit)
	} else if page > 0 && limit > 0 {
		q = q.Order("Timestamp").Offset((page - 1) * limit).Limit(limit)
	} else {
		// WARNING: this returns everything !
		q = q.Order("Timestamp")
	}

	_, err := q.GetAll(ctx, &events)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		events = make([]types.EventDS, 0)
	}

	logger.Info(ctx, topic, "ClientID=%s,time[%d,%d],page=%d,limit=%d. Found=%d", clientID, start, end, page, limit, len(events))
	return &events, nil
}

// StoreEvent stores an event in the datastore
func StoreEvent(ctx context.Context, clientID string, event *api.Event) error {
	topic := "backend.events.store"

	// deep copy of the struct
	e := types.EventDS{
		ClientID:         clientID,
		Event:            event.Event,
		EntityType:       event.EntityType,
		EntityID:         event.EntityID,
		TargetEntityType: event.TargetEntityType,
		TargetEntityID:   event.TargetEntityID,
		Properties:       event.Properties,
		Timestamp:        event.Timestamp,
		Created:          util.Timestamp(),
	}

	key := datastore.NewIncompleteKey(ctx, types.DatastoreEvents, nil)
	_, err := datastore.Put(ctx, key, &e)
	if err != nil {
		logger.Error(ctx, topic, err.Error())
	}

	return err
}

// ExportEvents exports events in time range ]start, end] and writes it to a csv file on Cloud Storage
func ExportEvents(ctx context.Context, exportID string) (int, error) {
	topic := "backend.events.export"
	set := make(map[string]bool)

	p := strings.Split(exportID, ".")
	clientID := p[0]
	event := p[1]

	export, err := GetExport(ctx, clientID, event)
	if err != nil {
		logger.Warning(ctx, topic, "Export not found. Export='%s'", exportID)
		return -1, err
	}

	// create a blob on Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Warning(ctx, topic, "Can not access storage. Export='%s'", exportID)
		return -1, err
	}

	// timerange for the export: ]start, end]
	start := export.LastExported
	end := util.Timestamp()
	numEvents := 0

	// monster query
	var q *datastore.Query
	if event == types.AllEvents {
		q = datastore.NewQuery(types.DatastoreEvents).Filter("ClientID =", clientID).Filter("Timestamp >", start).Limit(types.ExportBatchSize).Order("Timestamp")
	} else {
		q = datastore.NewQuery(types.DatastoreEvents).Filter("ClientID =", clientID).Filter("Event =", event).Filter("Timestamp >", start).Limit(types.ExportBatchSize).Order("Timestamp")
	}

	fileName := fmt.Sprintf("%s/parts/%s.%d.csv", clientID, event, start)
	bucket := client.Bucket(api.DefaultExportBucket)

	w := bucket.Object(fileName).NewWriter(ctx)
	w.ContentType = "text/plain"
	defer w.Close()

	// run the query and write the blob
	iter := q.Run(ctx)
	for {
		var e types.EventDS

		_, err := iter.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			logger.Warning(ctx, topic, "Could not query events. Export='%s'", exportID)
			return -1, err
		}

		w.Write([]byte(e.ToCSV()))
		set[e.Event] = true

		end = e.Timestamp
		numEvents++
	}

	if numEvents == 0 {
		// cleanup since nothing was exported
		w.Close()
		bucket.Object(fileName).Delete(ctx)
	}

	// create new exports if a new event type was found
	for e := range set {
		q = datastore.NewQuery(types.DatastoreExports).Filter("ClientID =", clientID).Filter("Event =", e)
		n, err := q.Count(ctx)
		if err != nil {
			logger.Warning(ctx, topic, "Could not query exports. Export='%s.%s'", clientID, e)
		}

		if n == 0 {
			_, err := CreateExport(ctx, clientID, e)
			if err != nil {
				logger.Warning(ctx, topic, "Could not create new export. Export='%s.%s'", clientID, e)
			} else {
				logger.Info(ctx, topic, "Created a new export. Export='%s.%s'", clientID, e)
			}
		}
	}

	logger.Info(ctx, topic, "Exported %d events. File='%s'", numEvents, fileName)

	// update metadata
	err = markExported(ctx, clientID, event, numEvents, end, util.IncT(util.Timestamp(), export.ExportSchedule))
	if err != nil {
		logger.Warning(ctx, topic, "Could not update metadata. Export='%s'", exportID)
		return -1, err
	}

	return numEvents, nil
}

// MergeEvents merges all exported events for a model in a single file
func MergeEvents(ctx context.Context, exportID string) error {
	var size int64
	var numFiles int

	topic := "backend.events.merge"

	p := strings.Split(exportID, ".")
	clientID := p[0]
	event := p[1]

	// get access to Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Warning(ctx, topic, "Could not access storage. Export='%s'", exportID)
		return err
	}

	// buckets
	sourceBucket := client.Bucket(api.DefaultExportBucket)
	targetBucket := client.Bucket(api.DefaultExportBucket)

	// new target blob
	fileName := fmt.Sprintf("%s/%s.csv", clientID, event)
	w := targetBucket.Object(fileName).NewWriter(ctx)
	w.ContentType = "text/plain"
	defer w.Close()

	// query blobs
	q := storage.Query{Prefix: clientID + "/parts/" + event}
	it := sourceBucket.Objects(ctx, &q)

	// merge the result
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. Export='%s'", exportID)
			return err
		}
		r, err := sourceBucket.Object(obj.Name).NewReader(ctx)
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. Export='%s'", exportID)
			return err
		}
		defer r.Close()

		// copy from one blob into the other
		b, err := io.Copy(w, r)
		if err != nil {
			logger.Warning(ctx, topic, "Could not copy exported events. Export='%s'", exportID)
			return err
		}

		numFiles++
		size += b
	}

	logger.Info(ctx, topic, "Merged %d files. Size=%d bytes", numFiles, size)
	return nil
}
