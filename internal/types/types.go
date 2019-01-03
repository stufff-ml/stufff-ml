package types

import "github.com/stufff-ml/stufff-ml/pkg/api"

const (
	// DatastoreEvents collection EVENTS
	DatastoreEvents string = "EVENTS"
	// DatastorePredictions collection PREDICTIONS
	DatastorePredictions string = "PREDICTIONS"
	// DatastoreModels collection MODELS
	DatastoreModels string = "MODELS"
	// DatastoreAuthorizations collection AUTHORIZATIONS
	DatastoreAuthorizations string = "AUTHORIZATIONS"
	// DatastoreClientResources collection CLIENT_RESOURCES
	DatastoreClientResources string = "CLIENT_RESOURCES"
	// DatastoreExports collection EXPORTS
	DatastoreExports string = "EXPORTS"

	// ShortCacheDuration default time to keep stuff in memory
	ShortCacheDuration string = "1m"
	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "10m"

	// BackgroundWorkQueue is the default background job queue
	BackgroundWorkQueue string = "background-work"

	// ScopeAdmin grants access to all operations
	ScopeAdmin string = "admin"

	// ExportBatchSize is the number of events to be exported in one job
	ExportBatchSize int = 10001
)

type (
	// EventDS adds metadata to the Event struct needed for internal purposes
	EventDS struct {
		ClientID         string   `json:"-"`
		Event            string   `json:"event"`
		EntityType       string   `json:"entity_type"`
		EntityID         string   `json:"entity_id"`
		TargetEntityType string   `json:"target_entity_type,omitempty"`
		TargetEntityID   string   `json:"target_entity_id,omitempty"`
		Properties       []string `datastore:",noindex" json:"properties,omitempty"`
		Timestamp        int64    `json:"timestamp"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// PredictionDS stores the materialized predictions for fast retrieval
	PredictionDS struct {
		ClientID string `json:"client_id"`
		Domain   string `json:"domain"`
		EntityID string `json:"entity_id"`
		Revision int    `json:"revision"`

		Items []api.ItemScore `datastore:",noindex" json:"items"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ModelDS represents a training model
	ModelDS struct {
		ClientID string `json:"client_id"`
		Name     string `json:"name"`
		Revision int    `json:"revision"`

		// Metadata
		TrainingSchedule int   `json:"training_schedule"`
		NextSchedule     int64 `json:"next_schedule"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ExportDS represents the export configuration of one clients data
	ExportDS struct {
		ClientID string `json:"client_id"`
		Event    string `json:"event"`

		// Metadata
		ExportSchedule int   `json:"export_schedule"`
		NextSchedule   int64 `json:"next_schedule"`
		LastExported   int64 `json:"exported"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ClientResourceDS represents an entity owning a client i.e. external source
	ClientResourceDS struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// AuthorizationDS represents access to a resource
	AuthorizationDS struct {
		ClientID string `json:"client_id"`
		Scope    string `json:"scope"`
		Token    string `json:"token"`
		Revoked  bool   `json:"revoked"`
		Expires  int64  `json:"expires"`

		// internal metadata
		Created int64 `json:"-"`
	}
)