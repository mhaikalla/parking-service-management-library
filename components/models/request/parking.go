package request

type ParkingInRequest struct {
	PlatNomor string `json:"plat_nomor" validate:"required"`
	Warna     string `json:"warna" validate:"required"`
	Tipe      string `json:"tipe" validate:"required"`
}

type ParkingOutRequest struct {
	PlatNomor string `json:"plat_nomor" validate:"required"`
}
