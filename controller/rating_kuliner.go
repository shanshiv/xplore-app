package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

// CreateRatingKuliner menambahkan rating kuliner baru
func CreateRatingKuliner(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingKuliner

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

	// Menyimpan data rating kuliner ke database
	if result := config.DB.Create(&rating); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Mengirimkan status Created (201) dan data rating yang baru disimpan
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

// GetAllRatingKuliner untuk mendapatkan semua rating kuliner
func GetAllRatingKuliner(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingKuliner
	if result := config.DB.Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// GetRatingKulinerByID untuk mendapatkan rating kuliner berdasarkan ID
func GetRatingKulinerByID(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingKuliner
	id := mux.Vars(r)["id"]

	// Menyaring rating berdasarkan id_kuliner
	if result := config.DB.Where("id_kuliner = ?", id).Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// UpdateRatingKuliner untuk memperbarui rating kuliner berdasarkan ID
func UpdateRatingKuliner(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingKuliner
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

// DeleteRatingKuliner untuk menghapus rating kuliner berdasarkan ID
func DeleteRatingKuliner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.RatingKuliner{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
