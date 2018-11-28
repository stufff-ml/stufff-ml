package backend

const (
	// DatastoreEvents collection EVENTS
	DatastoreEvents string = "EVENTS"
)

type (
	// EventsStore adds metadata to the Event struct needed for internal purposes
	EventsStore struct {
		AppDomain        string   `json:"-"`
		Event            string   `json:"event"`
		EntityType       string   `json:"entity_type"`
		EntityID         string   `json:"entity_id"`
		TargetEntityType string   `json:"target_entity_type,omitempty"`
		TargetEntityID   string   `json:"target_entity_id,omitempty"`
		Properties       []string `datastore:",noindex" json:"properties,omitempty"`
		Timestamp        int64    `json:"timestamp"`
		Created          int64    `json:"-"`
	}
)
