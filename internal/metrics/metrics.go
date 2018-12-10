package metrics

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
)

const (
	datastoreMetrics string = "METRICS"
)

type (

	// Metric is a generic data structure to store metrics
	Metric struct {
		Topic   string // describes the metric
		Label   string // additional context, e.g. an id, name/value pairs, comma separated labels etc
		Type    string // the type, e.g. count,
		Created int64
	}

	// Counter is a metric to collect integer values
	Counter struct {
		Metric
		Value int64
	}
)

// Count collects a numeric counter value
func Count(ctx context.Context, topic, label string, value int) {

	m := Counter{}
	m.Topic = topic
	m.Label = label
	m.Type = "count"
	m.Created = util.Timestamp()
	m.Value = int64(value)

	key := datastore.NewIncompleteKey(ctx, datastoreMetrics, nil)
	_, err := datastore.Put(ctx, key, &m)

	if err != nil {
		logger.Error(ctx, "metrics.count", err.Error())
	}
}
