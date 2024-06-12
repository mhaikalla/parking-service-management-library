package response

import "time"

type ParkingInResponse struct {
	PlatNomor    string    `json:"plat_nomor"`
	ParkingLot   string    `json:"parking_lot"`
	TanggalMasuk time.Time `json:"tanggal_masuk"`
}

type ParkingOutResponse struct {
	PlatNomor     string    `json:"plat_nomor"`
	JumlahBayar   string    `json:"jumlah_bayar"`
	TanggalMasuk  time.Time `json:"tanggal_masuk"`
	TanggalKeluar time.Time `json:"tanggal_keluar"`
}

type GetDataParkingResponse struct {
	PlatNomor []string `json:"plat_nomor"`
}

type GetCountParkingResponse struct {
	JumlahKendaraan int `json:"jumlah_kendaraan"`
}
