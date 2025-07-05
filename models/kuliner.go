package models

import (
	"time"

	"github.com/lib/pq"
)

type Kuliner struct {
	IDKuliner        uint           `gorm:"primaryKey;autoIncrement" json:"id_kuliner"`
	FotoKuliner      pq.StringArray `gorm:"type:text[]" json:"foto_kuliner"`
	NamaKuliner      string         `json:"nama_kuliner"`
	LokasiKuliner    string         `json:"lokasi_kuliner"`
	KoordinatKuliner string         `json:"koordinat_kuliner"`
	KontakKuliner    string         `json:"kontak_kuliner"`
	Deskripsi        string         `json:"deskripsi"`
	Fasilitas        pq.StringArray `gorm:"type:text[]" json:"fasilitas"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (Kuliner) TableName() string {
	return "kuliner"
}
