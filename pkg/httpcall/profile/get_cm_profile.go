package profile

import (
	"context"
	"net/http"

	endpointFn "parking-service/pkg/endpoint/functions"
	"parking-service/pkg/errs"
	"parking-service/pkg/httpc"
)

// ResponseCM ...
type ResponseAPI struct {
	Success bool        `json:"success"`
	Error   *errs.Errs  `json:"-"`
	Data    VehicleData `json:"data"`
}

// IVehicle ...
type IVehicle interface {
	GetVehicle(ctx context.Context, msisdn string) ResponseAPI
}

// Vehicle context hold configuration
type profileCM struct {
	client   *http.Client
	url      string
	method   string
	tokenGen func() (string, error)
}

// getVehicle ..
func (p *profileCM) GetVehicle(ctx context.Context, msisdn string) ResponseAPI {
	resp := ResponseAPI{}
	bodyResp := VehiclePayloadResp{}

	bodyFailed := VehicleFailedPayloadResp{}

	req := httpc.UpstreamsRequest{
		URL:          p.url,
		Method:       p.method,
		BearerFn:     p.tokenGen,
		Client:       p.client,
		BodySuccess:  &bodyResp,
		BodyFailed:   &bodyFailed,
		SuccessCodes: []int{200, 201},
	}

	res := req.RequestWithContext(ctx)
	if res.GetError() != nil {
		resp.Error = res.GetError()
		return resp
	}

	resp.Data = VehicleData{}
	resp.Success = true
	return resp
}

// VehicleWithMapConfig create Profile object using map as config, http client, and token generator as args
func VehicleMapConfig(config map[string]map[string]interface{}, client *http.Client, tokenGen func() (string, error)) (IVehicle, error) {

	parsedConfig, err := endpointFn.ParseMapConfig(config, "cm_profile", VehicleRouteNotFound)

	if err != nil {
		return nil, err
	}

	return &profileCM{
		client:   client,
		method:   parsedConfig.Method,
		url:      parsedConfig.URL,
		tokenGen: tokenGen,
	}, nil
}
