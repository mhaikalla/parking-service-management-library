package response

type GetDetailVehicleResponse struct {
	BaseResponse
	Name                string `json:"name"`
	Type                string `json:"type"`
	FirstHourPrice      int    `json:"first_hour_price"`
	PricePerHourPercent int    `json:"price_per_hour_percent"`
}

type GetVehiclesResponse struct {
	Data []GetDetailVehicleResponse `json:"Data"`
}
