package response

type ParkingInResponse struct {
	PlatNomor    string `json:"plat_nomor"`
	ParkingLot   string `json:"parking_lot"`
	TanggalMasuk string `json:"tanggal_masuk"`
}

type ParkingOutResponse struct {
	PlatNomor     string `json:"plat_nomor"`
	JumlahBayar   string `json:"jumlah_bayar"`
	TanggalMasuk  string `json:"tanggal_masuk"`
	TanggalKeluar string `json:"tanggal_keluar"`
}
