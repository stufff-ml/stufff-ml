package backend

import "github.com/stufff-ml/stufff-ml/pkg/types"

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

	// ShortCacheDuration default time to keep stuff in memory
	ShortCacheDuration string = "1m"
	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "10m"

	// ScopeAdmin grants access to all operations
	ScopeAdmin string = "admin"
)

type (
	// EventsStore adds metadata to the Event struct needed for internal purposes
	EventsStore struct {
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

	// PredictionStore stores the materialized predictions for fast retrieval
	PredictionStore struct {
		ClientID string `json:"client_id"`
		Domain   string `json:"domain"`
		EntityID string `json:"entity_id"`
		Revision int    `json:"revision"`

		Items []types.ItemScore `datastore:",noindex" json:"items"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// Model represents a training model
	Model struct {
		ClientID string `json:"client_id"`
		Domain   string `json:"domain"`
		Revision int    `json:"revision"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// ClientResource represents an entity owning a client i.e. external source
	ClientResource struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`

		// internal metadata
		Created int64 `json:"-"`
	}

	// Authorization represents access to a resource
	Authorization struct {
		ClientID string `json:"client_id"`
		Scope    string `json:"scope"`
		Token    string `json:"token"`
		Revoked  bool   `json:"revoked"`
		Expires  int64  `json:"expires"`

		// internal metadata
		Created int64 `json:"-"`
	}
)
