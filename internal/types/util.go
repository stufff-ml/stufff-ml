package types

import (
	"fmt"
	"strings"
)

// ToCSV creates a csv strin gfrom the struct
func (e *EventDS) ToCSV() string {
	if len(e.Properties) == 0 {
		return fmt.Sprintf("%s,%s,%s,%s,%s,%d,''\n", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%d,'%s'\n", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp, strings.Join(e.Properties, ","))
}
