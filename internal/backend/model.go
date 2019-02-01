package backend

import (
	"context"

	"google.golang.org/appengine/datastore"

	"github.com/majordomusio/commons/pkg/gae/logger"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/stufff-ml/stufff-ml/internal/types"
)

// CreateDefaultModel creates an initial model definition
func CreateDefaultModel(ctx context.Context, clientID string) (*types.ModelDS, error) {

	model := types.ModelDS{
		ClientID: clientID,
		Name:     types.Default,
		Revision: types.DefaultRevision,
		ConfigParams: []types.Parameters{
			{Key: "PythonModule", Value: "model.task"},
			{Key: "RuntimeVersion", Value: "1.12"},
			{Key: "PythonVersion", Value: "2.7"},
		},
		HyperParams: []types.Parameters{
			{Key: "weights", Value: "True"},
			{Key: "latent_factors", Value: "5"},
			{Key: "num_iters", Value: "20"},
			{Key: "regularization", Value: "0.07"},
			{Key: "unobs_weight", Value: "0.01"},
			{Key: "wt_type", Value: "0"},
			{Key: "feature_wt_factor", Value: "130.0"},
			{Key: "feature_wt_exp", Value: "0.08"},
		},
		Events:           []string{types.Default},
		Version:          0,
		TrainingSchedule: 180,
		NextSchedule:     0,
		Created:          util.Timestamp(),
	}

	key := ModelKey(ctx, clientID, types.Default)
	_, err := datastore.Put(ctx, key, &model)
	if err != nil {
		logger.Error(ctx, "backend.model.create", err.Error())
		return nil, err
	}

	return &model, nil
}

// GetModel returns the latest version of a model based on the clientID and name
func GetModel(ctx context.Context, clientID, name string) (*types.ModelDS, error) {
	var model types.ModelDS

	key := ModelKey(ctx, clientID, name)
	err := datastore.Get(ctx, key, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}

// MarkTrained writes an export record back to the datastore with updated metadata
func MarkTrained(ctx context.Context, clientID, name string, trained, next int64) error {

	model, err := GetModel(ctx, clientID, name)
	if err != nil {
		return err
	}

	model.LastTrained = trained
	model.NextSchedule = next

	_, err = datastore.Put(ctx, ModelKey(ctx, clientID, name), &model)
	if err != nil {
		return err
	}

	return err
}
