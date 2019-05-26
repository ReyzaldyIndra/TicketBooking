package v1

type Book struct {
	NamaTiket   string `json:"namaTiket"`
	HargaTiket  int    `json:"hargaTiket"`
	KelasTiket  string `json:"kelasTiket"`
	LokasiTiket string `json:"lokasiTiket"`
}
