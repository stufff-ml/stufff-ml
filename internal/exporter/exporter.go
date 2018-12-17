package exporter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/stufff-ml/stufff-ml/internal/backend"
)

// ExportEvents queries events from the datastore and writes a csv file
func ExportEvents(clientID string, start int64, endpoint, token, dataHome string) error {

	bearer := fmt.Sprintf("Bearer %s", token)
	uri := fmt.Sprintf("%s%s/events?id=%s&start=%d", endpoint, "types.InternalAPINamespace", clientID, start)
	filePathAndName := fmt.Sprintf("%s/%s_%d.csv", dataHome, clientID, start)

	// Create a new request using http
	req, err := http.NewRequest("GET", uri, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)
	check(err)
	defer resp.Body.Close()

	// decode the json payload
	events := make([]backend.EventsStore, 0)
	err = json.NewDecoder(resp.Body).Decode(&events)
	check(err)

	// open the file for writing
	f, err := os.Create(filePathAndName)
	check(err)
	defer f.Close()

	// loop over the array and write events to the file
	for i := range events {
		_, err := f.WriteString(eventToString(&events[i]) + "\n")
		check(err)
	}

	// makes sure everything is on disc
	err = f.Sync()
	check(err)

	return nil
}

func eventToString(e *backend.EventsStore) string {
	if len(e.Properties) == 0 {
		return fmt.Sprintf("%s,%s,%s,%s,%s,%d", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp)
	}
	return fmt.Sprintf("%s,%s,%s,%s,%s,%d,%s", e.Event, e.EntityType, e.EntityID, e.TargetEntityType, e.TargetEntityID, e.Timestamp, strings.Join(e.Properties, ","))
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
