package request

type CreateParkingLotRequest struct {
	Name  string `json:"name" validate:"required"`
	Floor string `json:"floor" validate:"required"`
}

type UpdateParkingLotRequest struct {
	Id    string `json:"parking_lot_id" validate:"required,numeric"`
	Name  string `json:"name" validate:"required"`
	Floor string `json:"floor" validate:"required"`
}

type DeleteParkingLotRequest struct {
	ParkingLotId string `json:"parking_lot_id" validate:"required,numeric"`
}

type GetDetailParkingLotRequest struct {
	ParkingLotId string `json:"parking_lot_id" validate:"required,numeric"`
}

type GetParkingLotRequest struct {
	BaseGetListParams
}
