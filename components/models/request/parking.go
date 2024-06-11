package request

type ParkingInRequest struct {
	PlatNomor string `json:"plat_nomor" validate:"required"`
	Warna     string `json:"warna" validate:"required"`
	Tipe      string `json:"tipe" validate:"required"`
}

type ParkingOutRequest struct {
	PlatNomor string `json:"plat_nomor" validate:"required"`
}

type GetParkingData struct {
	Warna string `json:"warna" validate:"required"`
	Tipe  string `json:"tipe"`
}

type GetCountParkingData struct {
	Tipe string `json:"tipe" validate:"required"`
}
