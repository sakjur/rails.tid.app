package trafikverket

import (
	"fmt"
	"os"
)

var apiKeyEnvironmentVariable = "RAILS_TRAFIKVERKET_SECRET"

func getAPIKey() (string, error) {
	key := os.Getenv(apiKeyEnvironmentVariable)
	if key == "" {
		return "", fmt.Errorf("API key not set")
	}

	return key, nil
}
