package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"

	"github.com/majordomusio/commons/pkg/gae/logger"

	"github.com/stufff-ml/stufff-ml/internal/types"
)

// ModelKey key on collection MODELS
func ModelKey(ctx context.Context, clientID, name string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreModels, clientID+"."+name, 0, nil)
}

// ExportKey key on collection EXPORTS
func ExportKey(ctx context.Context, clientID, event string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreExports, clientID+"."+event, 0, nil)
}

// ClientResourceKey key on collection CLIENT_RESOURCES
func ClientResourceKey(ctx context.Context, clientID string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreClientResources, clientID, 0, nil)
}

// AuthorizationKey key on collection AUTHORIZATIONS
func AuthorizationKey(ctx context.Context, token string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreAuthorizations, token, 0, nil)
}

// TrainingJobKey key on collection TRAINING_JOBS
func TrainingJobKey(ctx context.Context, k string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastoreTrainingJobs, k, 0, nil)
}

// PredictionKeyString returns the composite key string for a prediction
func PredictionKeyString(clientID, domain, entityID, revision string) string {
	return clientID + "." + domain + "." + entityID + "." + revision
}

// PredictionKey key on collection PREDICTIONS
func PredictionKey(ctx context.Context, k string) *datastore.Key {
	return datastore.NewKey(ctx, types.DatastorePredictions, k, 0, nil)
}

// ScheduleJob is a shorthand to create a background job
func ScheduleJob(ctx context.Context, queue, request string) error {
	t := taskqueue.NewPOSTTask(request, nil)
	_, err := taskqueue.Add(ctx, t, queue)
	if err != nil {
		logger.Error(ctx, "jobs.schedule", err.Error())
	}

	return err
}

// InvokeFunction calls a Cloud Function and posts data to it
func InvokeFunction(ctx context.Context, function, reqID string, payload interface{}) (int, *types.GenericResponse) {
	topic := "invoke.function"
	region := os.Getenv("REGION")
	projectID := os.Getenv("PROJECT_ID")
	uri := fmt.Sprintf("https://%s-%s.cloudfunctions.net/%s", region, projectID, function)

	r := types.GenericRequest{
		ReqID:   reqID,
		Payload: payload,
	}
	b, _ := json.Marshal(r)

	client := urlfetch.Client(ctx)
	client.Timeout = 55000

	req, _ := http.NewRequest("POST", uri, bytes.NewBuffer(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(ctx, topic, "ReqID=%s. Error=%s", reqID, err.Error())
		return http.StatusInternalServerError, nil
	}
	defer resp.Body.Close()

	var response types.GenericResponse
	json.NewDecoder(resp.Body).Decode(response)

	return resp.StatusCode, &response
}
