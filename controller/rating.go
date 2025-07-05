package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

// Fungsi untuk mengambil transaksi yang sudah diulas
func GetTransaksiWithReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"] // Ambil user_id dari URL params

	var transaksi []models.Transaksi
	if err := config.DB.Where("user_id = ?", userID).Find(&transaksi).Error; err != nil {
		http.Error(w, "Transaksi tidak ditemukan", http.StatusNotFound)
		return
	}

	var transaksiDetails []map[string]interface{}
	for _, t := range transaksi {
		data := make(map[string]interface{})
		data["id"] = t.ID
		data["user_id"] = t.UserID
		data["tipe_transaksi"] = t.TipeTransaksi
		data["total_harga"] = t.TotalHarga
		data["tipe_pembayaran"] = t.TipePembayaran
		data["created_at"] = t.CreatedAt

		// Ambil data berdasarkan tipe transaksi (kuliner, wisata, penginapan)
		switch t.TipeTransaksi {
		case "kuliner":
			var kuliner models.Kuliner
			if err := config.DB.First(&kuliner, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Kuliner tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek apakah sudah ada rating di tabel rating_kuliner
			var rating models.RatingKuliner
			if err := config.DB.Where("id_akun = ? AND id_kuliner = ?", t.UserID, t.IDTransaksi).First(&rating).Error; err == nil {
				// Rating ditemukan, berarti sudah diulas
				data["rating"] = rating.Rating
				data["ulasan"] = rating.Ulasan
				data["kuliner"] = kuliner
				transaksiDetails = append(transaksiDetails, data)
			}
		case "wisata":
			var wisata models.Wisata
			if err := config.DB.First(&wisata, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Wisata tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek apakah sudah ada rating di tabel rating_wisata
			var ratingWisata models.RatingWisata
			if err := config.DB.Where("id_akun = ? AND id_wisata = ?", t.UserID, t.IDTransaksi).First(&ratingWisata).Error; err == nil {
				// Rating ditemukan, berarti sudah diulas
				data["rating"] = ratingWisata.Rating
				data["ulasan"] = ratingWisata.Ulasan
				data["wisata"] = wisata
				transaksiDetails = append(transaksiDetails, data)
			}
		case "penginapan":
			var penginapan models.Penginapan
			if err := config.DB.First(&penginapan, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Penginapan tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek apakah sudah ada rating di tabel rating_penginapan
			var ratingPenginapan models.RatingPenginapan
			if err := config.DB.Where("id_akun = ? AND id_penginapan = ?", t.UserID, t.IDTransaksi).First(&ratingPenginapan).Error; err == nil {
				// Rating ditemukan, berarti sudah diulas
				data["rating"] = ratingPenginapan.Rating
				data["ulasan"] = ratingPenginapan.Ulasan
				data["penginapan"] = penginapan
				transaksiDetails = append(transaksiDetails, data)
			}
		default:
			http.Error(w, "Tipe transaksi tidak valid", http.StatusBadRequest)
			return
		}
	}

	// Kirimkan data transaksi yang sudah diulas ke client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaksiDetails)
}

// Fungsi untuk mengambil transaksi yang belum diulas
func GetTransaksiWithoutReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"] // Ambil user_id dari URL params

	var transaksi []models.Transaksi
	if err := config.DB.Where("user_id = ?", userID).Find(&transaksi).Error; err != nil {
		http.Error(w, "Transaksi tidak ditemukan", http.StatusNotFound)
		return
	}

	var transaksiDetails []map[string]interface{}
	for _, t := range transaksi {
		data := make(map[string]interface{})
		data["id"] = t.ID
		data["user_id"] = t.UserID
		data["tipe_transaksi"] = t.TipeTransaksi
		data["total_harga"] = t.TotalHarga
		data["tipe_pembayaran"] = t.TipePembayaran
		data["created_at"] = t.CreatedAt

		// Ambil data berdasarkan tipe transaksi (kuliner, wisata, penginapan)
		switch t.TipeTransaksi {
		case "kuliner":
			var kuliner models.Kuliner
			if err := config.DB.First(&kuliner, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Kuliner tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek jika belum ada rating di tabel rating_kuliner
			var rating models.RatingKuliner
			if err := config.DB.Where("id_akun = ? AND id_kuliner = ?", t.UserID, t.IDTransaksi).First(&rating).Error; err != nil {
				// Jika tidak ada rating, tampilkan transaksi
				data["kuliner"] = kuliner
				transaksiDetails = append(transaksiDetails, data)
			}
		case "wisata":
			var wisata models.Wisata
			if err := config.DB.First(&wisata, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Wisata tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek jika belum ada rating di tabel rating_wisata
			var ratingWisata models.RatingWisata
			if err := config.DB.Where("id_akun = ? AND id_wisata = ?", t.UserID, t.IDTransaksi).First(&ratingWisata).Error; err != nil {
				// Jika tidak ada rating, tampilkan transaksi
				data["wisata"] = wisata
				transaksiDetails = append(transaksiDetails, data)
			}
		case "penginapan":
			var penginapan models.Penginapan
			if err := config.DB.First(&penginapan, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Penginapan tidak ditemukan", http.StatusNotFound)
				return
			}

			// Cek jika belum ada rating di tabel rating_penginapan
			var ratingPenginapan models.RatingPenginapan
			if err := config.DB.Where("id_akun = ? AND id_penginapan = ?", t.UserID, t.IDTransaksi).First(&ratingPenginapan).Error; err != nil {
				// Jika tidak ada rating, tampilkan transaksi
				data["penginapan"] = penginapan
				transaksiDetails = append(transaksiDetails, data)
			}
		default:
			http.Error(w, "Tipe transaksi tidak valid", http.StatusBadRequest)
			return
		}
	}

	// Kirimkan data transaksi yang belum diulas ke client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaksiDetails)
}
