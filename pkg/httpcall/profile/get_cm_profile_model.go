package profile

const (
	VehicleRouteNotFound = "No configuration key vehicle found in upstreams_api entry"
)

type VehicleData struct {
}

type VehiclePayloadResp struct {
}

type VehicleFailedPayloadResp struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
