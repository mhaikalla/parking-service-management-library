package request

type CreateVehicleRequest struct {
	Name                string `json:"name" validate:"required"`
	Type                string `json:"type" validate:"required"`
	FirstHourPrice      int    `json:"first_hour_price" validate:"required"`
	PricePerHourPercent int    `json:"price_per_hour_percent" validate:"required"`
}

type UpdateVehicleRequest struct {
	Id                  int    `json:"vehicle_id" validate:"required"`
	Name                string `json:"name" validate:"required"`
	Type                string `json:"type" validate:"required"`
	FirstHourPrice      int    `json:"first_hour_price" validate:"required"`
	PricePerHourPercent int    `json:"price_per_hour_percent" validate:"required"`
}

type DeleteVehicleRequest struct {
	VehicleId string `json:"vehicle_id" validate:"required"`
}

type GetDetailVehicleRequest struct {
	VehicleId string `json:"vehicle_id" validate:"required"`
}

type GetVehicleRequest struct {
	BaseGetListParams
}
