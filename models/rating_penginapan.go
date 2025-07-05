package models

import (
	"time"
)

type RatingPenginapan struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IDPenginapan  uint      `gorm:"not null" json:"id_penginapan"`
	Rating        float64   `gorm:"not null" json:"rating"`
	Ulasan        string    `json:"ulasan"`
	TanggalUlasan time.Time `json:"tanggal_ulasan"`
	IDAkun        uint      `json:"id_akun"`  // Foreign key yang mengarah ke akun(id)
	Email         string    `json:"email"`    // Menyimpan email dari akun
	Username      string    `json:"username"` // Menyimpan username dari akun
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (RatingPenginapan) TableName() string {
	return "rating_penginapan"
}
