package models

const ParkingVehicleStatusTableName = "parking_vehicle_status"

type ParkingVehicleStatus struct {
	BaseEntity
	Name           string  `json:"name"`
	PlateNumber    string  `json:"plate_number"`
	Type           string  `json:"type"`
	Color          string  `json:"color"`
	ParkingInDate  string  `json:"parking_in_date"`
	ParkingOutDate *string `json:"parking_out_date"`
	Status         int     `json:"status"`
	Price          int     `json:"price"`
	ParkingLot     string  `json:"parking_lot"`
}
