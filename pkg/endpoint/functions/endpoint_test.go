package functions

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"parking-service/pkg/endpoint/models"
)

func TestParseMapConfig(t *testing.T) {
	type args struct {
		config        models.EndpointMapConfig
		keyConfig     string
		notFoundError string
	}
	tests := []struct {
		name    string
		args    args
		want    models.EndpointConfig
		wantErr bool
	}{
		{
			"upstreams entry not found",
			args{
				config:        map[string]map[string]interface{}{},
				keyConfig:     "test",
				notFoundError: models.UpstreamsConfigNotFound,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"route invalid type",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": nil,
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"route not found",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": {},
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"one of key on route not found",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": {
						"test": map[string]interface{}{
							"method": "POST",
						},
					},
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"one of key on route is empty string",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": {
						"test": map[string]interface{}{
							"method": "POST",
							"url":    "",
						},
					},
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"one of key on route is not string",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": {
						"test": map[string]interface{}{
							"method": 1,
							"url":    "/some/thing/foobar",
						},
					},
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{},
			true,
		},
		{
			"success",
			args{
				config: map[string]map[string]interface{}{
					"upstreams_api": {
						"test": map[string]interface{}{
							"method": "POST",
							"url":    "/some/thing/foobar",
						},
					},
				},
				keyConfig:     "test",
				notFoundError: models.IncompleteConfigMap,
			},
			models.EndpointConfig{
				Method: "POST",
				URL:    "/some/thing/foobar",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMapConfig(tt.args.config, tt.args.keyConfig, tt.args.notFoundError)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseMapConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseMapConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateHTTPClient(t *testing.T) {
	type args struct {
		config map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want *http.Client
	}{
		{
			"with conf",
			args{config: map[string]interface{}{
				"upstreams": map[string]interface{}{
					"timeout": 60,
				},
			}},
			&http.Client{Timeout: time.Second * 60},
		},
		{
			"without conf",
			args{config: map[string]interface{}{}},
			&http.Client{Timeout: time.Second * 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateHTTPClient(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateHTTPClient() = %v, want %v", got, tt.want)
			}
		})
	}
}
