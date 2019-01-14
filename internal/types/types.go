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

	// ScopeAdminFull grants access to all operations
	ScopeAdminFull string = "admin:full"
	// ScopeAPIFull grants access to all API operations
	ScopeAPIFull string = "api:full"
	// ScopeUserFull grants access to all API operations
	ScopeUserFull string = "user:full"

	// ScopeRootAccess gets you all access
	ScopeRootAccess string = "admin:full api:full"

	// ExportBatchSize is the number of events to be exported in one job
	ExportBatchSize int = 10001

	// Default is the key used to setup the default model etc
	Default string = "default"
	// AllEvents is the event type used to export everything
	AllEvents string = "all"
	// DefaultRevision is the default version number for models
	DefaultRevision int = 1
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
		Version  int    `json:"version"`

		Items []api.ItemScore `datastore:",noindex" json:"items"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ModelDS represents a training model
	ModelDS struct {
		ClientID string `json:"client_id"`
		Name     string `json:"name"`     // name of the model
		Revision int    `json:"revision"` // version of the model

		// Model configuration
		Version int      `json:"version"` // version of the model setup, params etc
		Events  []string `json:"events"`  // list of event types the model uses. Set to 'default'.

		// Model config params
		ConfigParams []Parameters `json:"config_params"`
		HyperParams  []Parameters `json:"hyper_params"`

		// Metadata
		TrainingSchedule int   `json:"training_schedule"`
		LastTrained      int64 `json:"last_trained"`
		NextSchedule     int64 `json:"next"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ExportDS represents the export configuration of one clients data
	ExportDS struct {
		ClientID string `json:"client_id"`
		Event    string `json:"event"`

		Exported      int `json:"exported"`
		ExportedTotal int `json:"exported_total"`

		// Metadata
		ExportSchedule int   `json:"export_schedule"`
		LastExported   int64 `json:"last_exported"`
		NextSchedule   int64 `json:"next"`

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

	// Parameters is a generic struct to store configuration parameters
	Parameters struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	GenericRequest struct {
		ReqID   string      `json:"id"`
		Payload interface{} `json:"payload"`
	}

	GenericResponse struct {
		Status  string `json:"status"` // ok | error
		Message string `json:"message,omitempty"`
	}
)
