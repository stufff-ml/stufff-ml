package types

type (

	// Event allows to reveive real-time event data.
	// See http://predictionio.apache.org/datacollection/eventmodel/
	Event struct {
		Event            string   `json:"event" binding:"required"`
		EntityType       string   `json:"entity_type" binding:"required"`
		EntityID         string   `json:"entity_id" binding:"required"`
		TargetEntityType string   `json:"target_entity_type,omitempty"`
		TargetEntityID   string   `json:"target_entity_id,omitempty"`
		Properties       []string `json:"properties,omitempty"`
		Timestamp        int64    `json:"timestamp,omitempty"`
	}

	// Prediction returns a set of predictions
	Prediction struct {
		EntityID string      `json:"entity_id" binding:"required"`
		Domain   string      `json:"domain" binding:"required"`
		Items    []ItemScore `json:"items"`
	}

	// ItemScore holds a single item recommendation and its score
	ItemScore struct {
		Item  string  `json:"item"`
		Score float64 `json:"score"`
	}

	// ClientResource returns a new client resource and its access token
	ClientResource struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Token        string `json:"token"`
	}
)
