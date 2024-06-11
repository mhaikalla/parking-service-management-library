package functions

import (
	"errors"
	"net/http"
	"time"

	"github.com/mhaikalla/parking-service-management-library/pkg/endpoint/models"
)

// ParseMapConfig parse map config into EndpointConfig object
func ParseMapConfig(config models.EndpointMapConfig, keyConfig, notFoundError string) (models.EndpointConfig, error) {
	res := models.EndpointConfig{}

	conf, confFound := config["upstreams_api"]

	if !confFound {
		return res, errors.New(models.UpstreamsConfigNotFound)
	}

	selectedConf, routeFound := conf[keyConfig]

	if !routeFound {
		return res, errors.New(notFoundError)
	}

	route, typeRoute := selectedConf.(map[string]interface{})

	if !typeRoute {
		return res, errors.New(notFoundError)
	}

	url, urlOk := route["url"].(string)
	method, methodOk := route["method"].(string)

	if !urlOk || !methodOk || url == "" || method == "" {
		return res, errors.New(notFoundError)
	}

	res.Method = method
	res.URL = url

	return res, nil
}

// CreateHTTPClient function to create HTTP client based on configuration provided
func CreateHTTPClient(config map[string]interface{}) *http.Client {
	if cconf, found := config["upstreams"].(map[string]interface{}); found {
		return &http.Client{Timeout: time.Second * time.Duration(cconf["timeout"].(int))}
	}
	return &http.Client{Timeout: time.Second * 30}
}
