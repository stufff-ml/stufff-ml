package backend

import (
	"context"
	"log"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/googleapi"
	ml "google.golang.org/api/ml/v1"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/types"
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
	// Google ML
	//

	ts, err := google.DefaultTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		log.Fatal(err)
	}
	client := oauth2.NewClient(ctx, ts)

	/*
		data, err := ioutil.ReadFile("../../test_user.json")
		if err != nil {
			log.Fatal(err)
		}

		creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/cloud-platform")
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Transport: &oauth2.Transport{
				Source: creds.TokenSource,
				//Source: google.AppEngineTokenSource(ctx, "https://www.googleapis.com/auth/cloud-platform"),
				Base: &urlfetch.Transport{
					Context: ctx,
				},
			},
		}
	*/

	//mlService, err := ml.New(urlfetch.Client(ctx))
	mlService, err := ml.New(client)
	if err != nil {
		logger.Warning(ctx, topic, "Could not create client for ML service. Model='%s'", modelID)
		return err
	}

	call := mlService.Projects.Jobs.List("projects/stufff-ml")
	resp, e := call.Do(googleapi.Trace("abcd1234"))
	//config := mlService.Projects.GetConfig("projects/stufff-ml")

	//getConfig := config.Context(ctx)
	//resp, e := getConfig.Do(googleapi.Trace("abcd1234"))

	logger.Info(ctx, topic, "+++ ", ts, resp, e)

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
