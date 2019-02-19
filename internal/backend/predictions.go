package backend

import (
	"bufio"
	"context"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/stufff-ml/stufff-ml/internal/types"
	"github.com/stufff-ml/stufff-ml/pkg/api"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine/datastore"
)

// GetPrediction returns a prediction based on a specified model
func GetPrediction(ctx context.Context, clientID string, req *api.Prediction) (*api.Prediction, error) {

	// lookup the prediction
	p := api.Prediction{}

	return &p, nil
}

// StorePrediction stores a materialized prediction in the datastore
func StorePrediction(ctx context.Context, prediction *api.Prediction) error {

	return nil
}

// ImportPredictions imports the results of a training job
func ImportPredictions(ctx context.Context, jobID string) error {
	var job types.TrainingJobDS
	topic := "backend.predictions.import"

	key := TrainingJobKey(ctx, jobID)
	err := datastore.Get(ctx, key, &job)
	if err != nil {
		logger.Warning(ctx, topic, "Could not load training data. JobID='%s'", jobID)
		return err
	}

	// get access to Cloud Storage
	client, err := storage.NewClient(ctx)
	if err != nil {
		logger.Warning(ctx, topic, "Could not access storage. JobID='%s'", jobID)
		return err
	}

	// query the bucket
	sourceBucket := client.Bucket(api.DefaultModelsBucket)
	q := storage.Query{Prefix: job.ClientID + "/" + jobID}
	it := sourceBucket.Objects(ctx, &q)

	// merge the result
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. JobID='%s'", jobID)
			return err
		}

		r, err := sourceBucket.Object(obj.Name).NewReader(ctx)
		if err != nil {
			logger.Warning(ctx, topic, "Could not access storage. JobID='%s'", jobID)
			return err
		}
		defer r.Close()

		n := 0
		s := bufio.NewScanner(r)
		for s.Scan() {
			pred := parse(s.Text())
			pred.ClientID = job.ClientID
			pred.ModelID = job.ModelID
			pred.Version = job.Version
			pred.Created = util.Timestamp()

			// FIXME batch store
			key := datastore.NewIncompleteKey(ctx, types.DatastorePredictions, nil)
			_, err := datastore.Put(ctx, key, pred)
			if err != nil {
				logger.Error(ctx, topic, err.Error())
			}
			n++
		}

		logger.Info(ctx, topic, "Imported %d predictions from '%s'. JobID='%s", n, obj.Name, jobID)
	}

	logger.Info(ctx, topic, "Importing of training job results finished. JobID=%s", jobID)

	return nil
}

// entity_id,entity_type,target_entity_type,values
// 1,user,item,"[7080, 0.99999, 968, 0.99999]"

func parse(p string) *api.Prediction {
	var pred api.Prediction

	parts := strings.Split(p, ",")

	// remove brackets etc
	parts[3] = parts[3][2:]
	last := len(parts) - 1
	parts[last] = parts[last][:len(parts[last])-2]

	pred.EntityID = parts[0]
	pred.EntityType = parts[1]
	pred.TargetEntityType = parts[2]

	values := parts[3:]
	t := len(values) / 2
	items := make([]api.ItemScore, t)

	for i := 0; i < t; i++ {
		items[i].EntityID = values[i*2]
		score, _ := strconv.ParseFloat(strings.Trim(values[(i*2)+1], " "), 64)
		items[i].Score = score
	}
	pred.Items = items

	return &pred
}
