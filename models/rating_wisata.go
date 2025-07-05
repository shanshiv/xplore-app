package models

import (
	"time"
)

type RatingWisata struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IDWisata      uint      `gorm:"not null" json:"id_wisata"`
	Rating        float64   `gorm:"not null" json:"rating"`
	Ulasan        string    `json:"ulasan"`
	TanggalUlasan time.Time `json:"tanggal_ulasan"`
	IDAkun        uint      `json:"id_akun"`  // Foreign key yang mengarah ke akun(id)
	Email         string    `json:"email"`    // Menyimpan email dari akun
	Username      string    `json:"username"` // Menyimpan username dari akun
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (RatingWisata) TableName() string {
	return "rating_wisata"
}
