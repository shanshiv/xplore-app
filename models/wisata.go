package models

import (
	"time"

	"github.com/lib/pq"
)

type Wisata struct {
	IDWisata        uint           `gorm:"primaryKey;autoIncrement" json:"id_wisata"`
	FotoWisata      pq.StringArray `gorm:"type:text[]" json:"foto_wisata"`
	NamaWisata      string         `json:"nama_wisata"`
	LokasiWisata    string         `json:"lokasi_wisata"`
	KoordinatWisata string         `json:"koordinat_wisata"`
	KontakWisata    string         `json:"kontak_wisata"`
	Deskripsi       string         `json:"deskripsi"`
	Fasilitas       pq.StringArray `gorm:"type:text[]" json:"fasilitas"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (Wisata) TableName() string {
	return "wisata"
}
