package models

const VehicleTableName = "vehicle"

type Vehicle struct {
	BaseEntity
	Name                string `json:"name"`
	Type                string `json:"type"`
	FirstHourPrice      int    `json:"first_hour_price"`
	PricePerHourPercent int    `json:"price_per_hour_percent"`
}
