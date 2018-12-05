package main

import (
	"flag"
	"os"

	"github.com/stufff-ml/stufff-ml/internal/exporter"
	"github.com/stufff-ml/stufff-ml/pkg/types"
)

// go run exporter.go -token=xoxo-aaaaaaaa -id=foo1233

// CLIENT_TOKEN	-> access token of the caller
// ID						-> client_id for the data to be exported
// START				-> timestamp to start exporting from
// API_ENDPOINT	-> url of the API
// DATA_HOME		-> absolute path to the location of the files

func main() {

	// get the data from the command line
	clientID := flag.String("id", "", "ClientID of the data to be exported")
	token := flag.String("token", "", "Token of the caller process")
	endpoint := flag.String("url", types.DefaultEndpoint, "API endpoint url")
	dataHome := flag.String("data", ".", "Location of the data")
	start := flag.Int64("start", 0, "Timestamp of the oldest event")
	flag.Parse()

	err := exporter.ExportEvents(*clientID, *start, *endpoint, *token, *dataHome)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
