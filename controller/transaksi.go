package controller

import (
	"encoding/json"
	"net/http"
	"time"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

// Fungsi untuk membuat transaksi
func CreateTransaksi(w http.ResponseWriter, r *http.Request) {
	var transaksi models.Transaksi

	// Decode data JSON yang dikirimkan
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&transaksi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cek apakah id_transaksi ada dan tipe_transaksi valid
	if transaksi.TipeTransaksi == "" || transaksi.IDTransaksi == 0 {
		http.Error(w, "Tipe transaksi dan ID transaksi harus diisi", http.StatusBadRequest)
		return
	}

	// Mengecek ID transaksi berdasarkan tipe transaksi (kuliner, wisata, penginapan)
	var isValid bool
	switch transaksi.TipeTransaksi {
	case "kuliner":
		var kuliner models.Kuliner
		if err := config.DB.First(&kuliner, transaksi.IDTransaksi).Error; err != nil {
			http.Error(w, "Kuliner tidak ditemukan", http.StatusNotFound)
			return
		}
		isValid = true
	case "wisata":
		var wisata models.Wisata
		if err := config.DB.First(&wisata, transaksi.IDTransaksi).Error; err != nil {
			http.Error(w, "Wisata tidak ditemukan", http.StatusNotFound)
			return
		}
		isValid = true
	case "penginapan":
		var penginapan models.Penginapan
		if err := config.DB.First(&penginapan, transaksi.IDTransaksi).Error; err != nil {
			http.Error(w, "Penginapan tidak ditemukan", http.StatusNotFound)
			return
		}
		isValid = true
	default:
		http.Error(w, "Tipe transaksi tidak valid", http.StatusBadRequest)
		return
	}

	// Jika ID transaksi ditemukan, simpan transaksi ke database
	if isValid {
		transaksi.CreatedAt = time.Now()

		// Simpan transaksi ke database
		if err := config.DB.Create(&transaksi).Error; err != nil {
			http.Error(w, "Unable to save to database", http.StatusInternalServerError)
			return
		}

		// Kirim respons sukses
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transaksi)
	}
}

// Fungsi untuk mengambil semua transaksi berdasarkan user_id dan mengecek tipe transaksi
func GetTransaksiByUserID(w http.ResponseWriter, r *http.Request) {
	// Ambil ID User dari URL params
	vars := mux.Vars(r)
	userID := vars["id"]

	// Ambil data transaksi berdasarkan user_id
	var transaksi []models.Transaksi
	if err := config.DB.Where("user_id = ?", userID).Find(&transaksi).Error; err != nil {
		http.Error(w, "Transaksi tidak ditemukan", http.StatusNotFound)
		return
	}

	// Inisialisasi slice untuk menyimpan data transaksi yang sudah lengkap
	var transaksiDetails []map[string]interface{}

	// Loop melalui setiap transaksi untuk mencari tipe transaksi dan ID yang relevan
	for _, t := range transaksi {
		// Inisialisasi map untuk data transaksi
		data := make(map[string]interface{})
		data["id"] = t.ID
		data["user_id"] = t.UserID
		data["tipe_transaksi"] = t.TipeTransaksi
		data["total_harga"] = t.TotalHarga
		data["tipe_pembayaran"] = t.TipePembayaran
		data["created_at"] = t.CreatedAt

		// Ambil detail berdasarkan tipe transaksi
		switch t.TipeTransaksi {
		case "kuliner":
			var kuliner models.Kuliner
			if err := config.DB.First(&kuliner, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Kuliner tidak ditemukan", http.StatusNotFound)
				return
			}
			data["kuliner"] = kuliner
		case "wisata":
			var wisata models.Wisata
			if err := config.DB.First(&wisata, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Wisata tidak ditemukan", http.StatusNotFound)
				return
			}
			data["wisata"] = wisata
		case "penginapan":
			var penginapan models.Penginapan
			if err := config.DB.First(&penginapan, t.IDTransaksi).Error; err != nil {
				http.Error(w, "Penginapan tidak ditemukan", http.StatusNotFound)
				return
			}
			data["penginapan"] = penginapan
		default:
			http.Error(w, "Tipe transaksi tidak valid", http.StatusBadRequest)
			return
		}

		// Menambahkan data transaksi ke dalam array
		transaksiDetails = append(transaksiDetails, data)
	}

	// Kirimkan data transaksi yang sudah lengkap ke client
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(transaksiDetails)
}
