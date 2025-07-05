package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

// CreateRatingPenginapan menambahkan rating penginapan baru
func CreateRatingPenginapan(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingPenginapan

	// Mendekode JSON request body
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rating); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mencari akun berdasarkan ID akun yang diberikan
	var akun models.Akun
	if err := config.DB.Where("id = ?", rating.IDAkun).First(&akun).Error; err != nil {
		http.Error(w, "Akun tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menyisipkan email dan username dari akun ke dalam rating
	rating.Email = akun.Email
	rating.Username = akun.Username

	// Menyimpan data rating penginapan ke database
	if result := config.DB.Create(&rating); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Mengirimkan status Created (201) dan data rating yang baru disimpan
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

// GetAllRatingPenginapan untuk mendapatkan semua rating penginapan
func GetAllRatingPenginapan(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingPenginapan
	if result := config.DB.Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// GetRatingPenginapanByID untuk mendapatkan rating penginapan berdasarkan ID
func GetRatingPenginapanByID(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingPenginapan
	id := mux.Vars(r)["id"]

	// Menyaring rating berdasarkan id_penginapan
	if result := config.DB.Where("id_penginapan = ?", id).Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// UpdateRatingPenginapan untuk memperbarui rating penginapan berdasarkan ID
func UpdateRatingPenginapan(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingPenginapan
	id := mux.Vars(r)["id"]

	// Cek apakah rating ada
	if result := config.DB.First(&rating, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	// Decode data JSON dan update kolom
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rating); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := config.DB.Save(&rating); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rating)
}

// DeleteRatingPenginapan untuk menghapus rating penginapan berdasarkan ID
func DeleteRatingPenginapan(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.RatingPenginapan{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
