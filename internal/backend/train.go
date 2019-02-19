package backend

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/majordomusio/commons/pkg/errors"
	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine/datastore"

	"github.com/stufff-ml/stufff-ml/internal/types"
	"github.com/stufff-ml/stufff-ml/pkg/api"
)

// TrainingInput is used to submit a training job with Cloud ML
type TrainingInput struct {
	// FIXME: make this a generic struct for a generic glue function!
	ScaleTier      string   `json:"scaleTier"`
	PackageURIs    []string `json:"packageUris"`
	PythonModule   string   `json:"pythonModule"`
	Region         string   `json:"region"`
	JobDir         string   `json:"jobDir"`
	RuntimeVersion string   `json:"runtimeVersion"`
	PythonVersion  string   `json:"pythonVersion"`
	ModelArguments []string `json:"args"`
}

// TrainModel submits a model for training to Google ML
func TrainModel(ctx context.Context, modelID string) error {
	topic := "backend.model.train"

	p := strings.Split(modelID, ".")
	clientID := p[0]
	name := p[1]

	model, err := GetModel(ctx, clientID, name)
	if err != nil {
		logger.Warning(ctx, topic, "Model not found. Model='%s'", modelID)
		return err
	}

	// update the model before building the training request
	model.Version++
	_, err = datastore.Put(ctx, ModelKey(ctx, clientID, name), model)
	if err != nil {
		logger.Warning(ctx, topic, "Could not update the model. Model='%s'", modelID)
		return err
	}

	// data for the training job
	region := os.Getenv("REGION")
	jobID := fmt.Sprintf("%s_%s_%d", model.Name, model.ClientID, util.Timestamp())
	jobDir := fmt.Sprintf("gs://%s/%s/%s", api.DefaultModelsBucket, model.ClientID, jobID)
	packageName := fmt.Sprintf("%s-%d", model.Name, model.Revision)
	callback := fmt.Sprintf("%s/%s/train?id=%s&job=%s", api.APIBaseURL, api.CallbackPrefix, model.ClientID, jobID)
	uris := []string{fmt.Sprintf("gs://%s/packages/%s/%s.tar.gz", api.DefaultResourcesBucket, packageName, packageName)}
	args := []string{"--client-id", model.ClientID, "--model-name", model.Name, "--job-id", jobID, "--callback", callback}

	job := types.TrainingJobDS{
		ClientID:       clientID,
		ModelID:        fmt.Sprintf("%s.%s", model.ClientID, model.Name),
		Version:        model.Version,
		JobID:          jobID,
		JobStarted:     util.Timestamp(),
		ModelArguments: args,
		Status:         "undefined",
		Created:        util.Timestamp(),
	}

	key := TrainingJobKey(ctx, jobID)
	_, err = datastore.Put(ctx, key, &job)
	if err != nil {
		logger.Warning(ctx, topic, "Could not create training job. Job ID='%s'", modelID)
		return err
	}

	//
	// Invoke Google ML API glue code
	//

	trainingInput := TrainingInput{
		ScaleTier:      "BASIC",
		PackageURIs:    uris,
		PythonModule:   "model.train",
		Region:         region,
		JobDir:         jobDir,
		RuntimeVersion: "1.12",
		PythonVersion:  "2.7",
		ModelArguments: args,
	}

	status, _ := InvokeFunction(ctx, "train_model", jobID, &trainingInput)

	if status == http.StatusOK {
		err = MarkTrained(ctx, clientID, name, util.Timestamp(), util.IncT(util.Timestamp(), model.TrainingSchedule))
		if err != nil {
			logger.Warning(ctx, topic, "Could not update metadata. Model='%s'", modelID)
			return err
		}
		logger.Info(ctx, topic, "Submitted model %s.%s for training. Job ID='%s'", clientID, name, jobID)
	} else {
		logger.Warning(ctx, topic, "Error submitting model %s.%s for training. Job ID='%s'", clientID, name, jobID)
		return errors.New("Error submitting model")
	}

	return nil
}

// MarkModelTrainingDone writes an export record back to the datastore with updated metadata
func MarkModelTrainingDone(ctx context.Context, jobID, status string) error {
	var job types.TrainingJobDS
	topic := "backend.model.training.done"

	key := TrainingJobKey(ctx, jobID)
	err := datastore.Get(ctx, key, &job)
	if err != nil {
		logger.Warning(ctx, topic, "Could not load training data. JobID='%s'", jobID)
		return err
	}

	// update the job record first
	job.JobEnded = util.Timestamp()
	job.Duration = job.JobEnded - job.JobStarted
	job.Status = status

	_, err = datastore.Put(ctx, key, &job)
	if err != nil {
		logger.Warning(ctx, topic, "Could not update training data. JobID='%s'", jobID)
	}

	return err
}
