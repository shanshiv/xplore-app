package models

import (
	"time"
)

type RatingKuliner struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IDKuliner     uint      `gorm:"not null" json:"id_kuliner"`
	Rating        float64   `gorm:"not null" json:"rating"`
	Ulasan        string    `json:"ulasan"`
	TanggalUlasan time.Time `json:"tanggal_ulasan"`
	IDAkun        uint      `json:"id_akun"`  // Foreign key yang mengarah ke akun(id)
	Email         string    `json:"email"`    // Menyimpan email dari akun
	Username      string    `json:"username"` // Menyimpan username dari akun
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (RatingKuliner) TableName() string {
	return "rating_kuliner"
}
