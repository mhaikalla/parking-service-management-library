package models

import "net/http"

const (
	// UpstreamsConfigNotFound error string if upstreams_api key not found in configuration
	UpstreamsConfigNotFound = "No configuration key upstreams_api found"

	// IncompleteConfigMap error string if config is incomplete
	IncompleteConfigMap = "Incomplete config map"
)

// EndpointConfig hold key pair configuration used on endpoint upstreams
type EndpointConfig struct {
	Method string `json:"method"`
	URL    string `json:"url"`
}

// EndpointContext ...
type EndpointContext struct {
	Config   EndpointConfig
	TokenGen func() (string, error)
	Client   *http.Client
}

// EndpointMapConfig map config for endpoint
type EndpointMapConfig map[string]map[string]interface{}

// EndpointConfigFn ...
type EndpointConfigFn func(structure interface{}) error
