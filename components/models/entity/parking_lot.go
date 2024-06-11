package models

const ParkingLotTableName = "parking_lot"

type ParkingLot struct {
	BaseEntity
	Name     string `json:"name"`
	Floor    string `json:"floor"`
	IsParked bool   `json:"isParked"`
}
