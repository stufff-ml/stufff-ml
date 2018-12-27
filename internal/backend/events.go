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

	"github.com/stufff-ml/stufff-ml/pkg/types"
)

// GetEvents queries the events store for events of type event in the time range [start, end]
func GetEvents(ctx context.Context, clientID, event string, start, end int64, page, limit int) (*[]EventDS, error) {
	topic := "backend.events.get"

	var events []EventDS
	var q *datastore.Query

	q = datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID)

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
		events = make([]EventDS, 0)
	}

	logger.Info(ctx, topic, "ClientID=%s,time[%d,%d],page=%d,limit=%d. Found=%d", clientID, start, end, page, limit, len(events))
	return &events, nil
}

// StoreEvent stores an event in the datastore
func StoreEvent(ctx context.Context, clientID string, event *types.Event) error {
	topic := "backend.events.store"

	// deep copy of the struct
	e := EventDS{
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
		logger.Error(ctx, topic, err.Error())
	}

	return err
}

// ExportEvents exports events in time range ]start, end] and writes it to a csv file on Cloud Storage
func ExportEvents(ctx context.Context, modelID string) (int, error) {
	topic := "backend.events.export"

	p := strings.Split(modelID, ".")
	clientID := p[0]
	domain := p[1]

	model, err := GetModel(ctx, clientID, domain)
	if err != nil {
		logger.Warning(ctx, topic, "Model not found. Model='%s'", modelID)
		return -1, err
	}

	// create a blob on Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Warning(ctx, topic, "Can not access storage. Model='%s'", modelID)
		return -1, err
	}

	// timerange for the export: ]start, end]
	start := model.LastExported
	end := util.Timestamp()
	numEvents := 0

	// monster query
	q := datastore.NewQuery(DatastoreEvents).Filter("ClientID =", clientID).Filter("Timestamp >", start).Limit(ExportBatchSize).Order("Timestamp")

	fileName := fmt.Sprintf("%s/%s.%d.csv", modelID, modelID, start)
	bucket := client.Bucket("exports.stufff.review")

	w := bucket.Object(fileName).NewWriter(ctx)
	w.ContentType = "text/plain"
	defer w.Close()

	// run the query and write the blob
	iter := q.Run(ctx)
	for {
		var e EventDS

		_, err := iter.Next(&e)
		if err == datastore.Done {
			break
		}
		if err != nil {
			logger.Warning(ctx, topic, "Could not query events. Model='%s'", modelID)
			return -1, err
		}

		w.Write([]byte(e.ToCSV()))

		end = e.Timestamp
		numEvents++
	}

	if numEvents == 0 {
		// cleanup since nothing was exported
		w.Close()
		bucket.Object(fileName).Delete(ctx)
	}

	logger.Info(ctx, topic, "Exported %d events. File='%s'", numEvents, fileName)

	// uodate metadata
	model.LastExported = end
	model.NextSchedule = util.IncT(end, model.TrainingSchedule)
	err = MarkModelExported(ctx, clientID, domain, end, util.IncT(util.Timestamp(), model.ExportSchedule))
	if err != nil {
		logger.Warning(ctx, topic, "Could not update metadata. Model='%s'", modelID)
		return -1, err
	}

	return numEvents, nil
}

// MergeEvents merges all exported events for a model in a single file
func MergeEvents(ctx context.Context, modelID string) error {
	var size int64
	var numFiles int

	topic := "backend.events.merge"

	p := strings.Split(modelID, ".")
	clientID := p[0]
	domain := p[1]

	model, err := GetModel(ctx, clientID, domain)
	if err != nil {
		logger.Warning(ctx, topic, "Model not found. Model='%s'", modelID)
		return err
	}

	// get access to Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Warning(ctx, topic, "Could not access storage. Model='%s'", modelID)
		return err
	}

	// buckets
	sourceBucket := client.Bucket(types.DefaultExportBucket)
	targetBucket := client.Bucket(types.DefaultModelsBucket)

	// new target blob
	fileName := fmt.Sprintf("%s/%s.%d.csv", modelID, modelID, model.Revision)
	w := targetBucket.Object(fileName).NewWriter(ctx)
	w.ContentType = "text/plain"
	defer w.Close()

	// query blobs
	q := storage.Query{Prefix: modelID}
	it := sourceBucket.Objects(ctx, &q)

	// merge the result
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. Model='%s'", modelID)
			return err
		}
		r, err := sourceBucket.Object(obj.Name).NewReader(ctx)
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. Model='%s'", modelID)
			return err
		}
		defer r.Close()

		// copy from one blob into the other
		b, err := io.Copy(w, r)
		if err != nil {
			logger.Warning(ctx, topic, "Could not copy exported events. Model='%s'", modelID)
			return err
		}

		numFiles++
		size += b
	}

	logger.Info(ctx, topic, "Merged %d files. Size=%d bytes", numFiles, size)
	return nil
}
