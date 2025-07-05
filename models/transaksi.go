package models

import (
	"time"
)

type Transaksi struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"not null" json:"user_id"`        // Foreign Key yang mengarah ke akun(id)
	TipeTransaksi  string    `gorm:"not null" json:"tipe_transaksi"` // kuliner, wisata, penginapan
	IDTransaksi    uint      `json:"id_transaksi"`                   // ID dari tabel kuliner, wisata, atau penginapan
	TotalHarga     float64   `json:"total_harga"`
	TipePembayaran string    `json:"tipe_pembayaran"`
	CreatedAt      time.Time `json:"created_at"` // Menyimpan waktu transaksi
}

func (Transaksi) TableName() string {
	return "transaksi"
}
