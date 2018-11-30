package backend

const (
	// DatastoreEvents collection EVENTS
	DatastoreEvents string = "EVENTS"
	// DatastoreModels collection MODELS
	DatastoreModels string = "MODELS"
	// DatastoreAuthorizations collection AUTHORIZATIONS
	DatastoreAuthorizations string = "AUTHORIZATIONS"
	// DatastoreClientResources collection CLIENT_RESOURCES
	DatastoreClientResources string = "CLIENT_RESOURCES"

	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "10m"
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

	// Model represents a training model
	Model struct {
		ModelID  string `json:"model_id"`
		ClientID string `json:"client_id"`
		Revision int    `json:"revision"`
		Event    string `json:"event"`

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
		Token    string `json:"token"`
		Revoked  bool   `json:"revoked"`
		Expires  int64  `json:"expires"`

		// internal metadata
		Created int64 `json:"-"`
	}
)
