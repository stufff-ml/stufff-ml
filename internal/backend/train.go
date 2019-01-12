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

	//
	// Invoke Google ML API glue code
	//

	region := os.Getenv("REGION")
	jobID := fmt.Sprintf("%s_%s_%d", model.Name, model.ClientID, util.Timestamp())
	jobDir := fmt.Sprintf("gs://%s/%s/%s", api.DefaultModelsBucket, model.ClientID, jobID)
	modelPackage := fmt.Sprintf("%s-%d", model.Name, model.Revision)
	uris := []string{fmt.Sprintf("gs://%s/packages/%s/%s.tar.gz", api.DefaultResourcesBucket, modelPackage, modelPackage)}
	args := []string{"--client-id", model.ClientID, "--model-name", model.Name, "--job-id", jobID}

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
		err = markTrained(ctx, clientID, name, util.Timestamp(), util.IncT(util.Timestamp(), model.TrainingSchedule))
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
