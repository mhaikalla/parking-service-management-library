package response

type GetDetailParkingLotResponse struct {
	BaseResponse
	Name     string `json:"name"`
	Floor    string `json:"floor"`
	IsParked bool   `json:"isParked"`
}

type GetParkingLotsResponse struct {
	Data []GetDetailParkingLotResponse `json:"data"`
}
