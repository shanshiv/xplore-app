package models

import (
	"time"

	"github.com/lib/pq"
)

type Penginapan struct {
	IDPenginapan        uint           `gorm:"primaryKey;autoIncrement" json:"id_penginapan"`
	FotoPenginapan      pq.StringArray `gorm:"type:text[]" json:"foto_penginapan"`
	NamaPenginapan      string         `json:"nama_penginapan"`
	LokasiPenginapan    string         `json:"lokasi_penginapan"`
	KoordinatPenginapan string         `json:"koordinat_penginapan"`
	KontakPenginapan    string         `json:"kontak_penginapan"`
	Deskripsi           string         `json:"deskripsi"`
	Fasilitas           pq.StringArray `gorm:"type:text[]" json:"fasilitas"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
}

func (Penginapan) TableName() string {
	return "penginapan"
}
