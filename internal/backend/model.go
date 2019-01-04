package backend

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/urlfetch"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/types"
	"github.com/stufff-ml/stufff-ml/pkg/api"
)

// CreateModel creates an initial model definition
func CreateModel(ctx context.Context, clientID, name string) (*types.ModelDS, error) {

	model := types.ModelDS{
		ClientID:         clientID,
		Name:             name,
		Revision:         1,
		TrainingSchedule: 60,
		NextSchedule:     0,
		Created:          util.Timestamp(),
	}

	key := ModelKey(ctx, clientID, name)
	_, err := datastore.Put(ctx, key, &model)
	if err != nil {
		logger.Error(ctx, "backend.model.create", err.Error())
		return nil, err
	}

	return &model, nil

}

// GetModel returns a model based on the clientID and domain
func GetModel(ctx context.Context, clientID, name string) (*types.ModelDS, error) {
	var model types.ModelDS

	// lookup the model definition
	key := "model." + strings.ToLower(clientID+"."+name)
	_, err := memcache.Gob.Get(ctx, key, &model)

	if err != nil {
		var models []types.ModelDS
		q := datastore.NewQuery(types.DatastoreModels).Filter("ClientID =", clientID).Filter("Name =", name).Order("-Revision")
		_, err := q.GetAll(ctx, &models)
		if err != nil {
			return nil, err
		}

		if len(models) == 0 {
			return nil, err
		}

		model = models[0]
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = model
			cache.Expiration, _ = time.ParseDuration(types.ShortCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &model, nil
}

// SubmitModel submits a model for training to Google ML
func SubmitModel(ctx context.Context, modelID string) error {
	topic := "backend.model.train"

	p := strings.Split(modelID, ".")
	clientID := p[0]
	name := p[1]

	model, err := GetModel(ctx, clientID, name)
	if err != nil {
		logger.Warning(ctx, topic, "Model not found. Model='%s'", modelID)
		return err
	}

	//
	// Invoke Google ML API glue code
	//

	modelPackage := fmt.Sprintf("%s-%d", model.Name, model.Revision)
	jobID := fmt.Sprintf("%s_%s_%d", model.Name, model.ClientID, util.Timestamp())
	jobDir := fmt.Sprintf("gs://%s/%s/%s", api.DefaultModelsBucket, model.ClientID, jobID)
	uris := []string{fmt.Sprintf("gs://models.stufff.review/packages/%s/%s.tar.gz", modelPackage, modelPackage)}
	args := []string{"--model-id", model.ClientID, "--model-event", model.Name}
	trainingInput := types.TrainingInput{
		ProjectID:      "stufff-review",
		JobID:          jobID,
		ScaleTier:      "BASIC",
		PackageURIs:    uris,
		PythonModule:   "model.task",
		Region:         "europe-west1",
		JobDir:         jobDir,
		RuntimeVersion: "1.12",
		PythonVersion:  "2.7",
		ModelArguments: args,
	}
	b, _ := json.Marshal(trainingInput)

	client := urlfetch.Client(ctx)
	client.Timeout = 55000
	req, _ := http.NewRequest("POST", "https://europe-west1-stufff-review.cloudfunctions.net/func_submit", bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Warning(ctx, topic, "Error submitting job. Model='%s'. Error=%s", modelID, err.Error())
		return err
	}
	defer resp.Body.Close()

	logger.Info(ctx, topic, "+++ DEBUG %d", resp.StatusCode)

	// update metadata
	err = markTrained(ctx, clientID, name, 0, util.IncT(util.Timestamp(), model.TrainingSchedule))
	if err != nil {
		logger.Warning(ctx, topic, "Could not update metadata. Model='%s'", modelID)
		return err
	}

	logger.Info(ctx, topic, "Submitted model %s.%s for training. Client='%s'", clientID, name, clientID)

	return nil
}

// MarkTrained writes an export record back to the datastore with updated metadata
func markTrained(ctx context.Context, clientID, name string, trained, next int64) error {
	var model types.ModelDS

	key := ModelKey(ctx, clientID, name)
	err := datastore.Get(ctx, key, &model)
	if err != nil {
		return err
	}

	//model.LastTrained = trained
	model.NextSchedule = next

	_, err = datastore.Put(ctx, key, &model)
	if err != nil {
		return err
	}

	// invalidate the cache
	ckey := "model." + strings.ToLower(clientID+"."+name)
	err = memcache.Delete(ctx, ckey)

	return err
}
