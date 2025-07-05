package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

func CreateRatingWisata(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingWisata

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

	// Menyimpan data rating wisata ke database
	if result := config.DB.Create(&rating); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Mengirimkan status Created (201) dan data rating yang baru disimpan
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rating)
}

// GetAllRatingWisata untuk mendapatkan semua rating wisata
func GetAllRatingWisata(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingWisata
	if result := config.DB.Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// GetRatingWisataByID untuk mendapatkan rating wisata berdasarkan ID
func GetRatingWisataByID(w http.ResponseWriter, r *http.Request) {
	var ratings []models.RatingWisata
	id := mux.Vars(r)["id"]

	// Menyaring rating berdasarkan id_wisata
	if result := config.DB.Where("id_wisata = ?", id).Find(&ratings); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ratings)
}

// UpdateRatingWisata untuk memperbarui rating wisata berdasarkan ID
func UpdateRatingWisata(w http.ResponseWriter, r *http.Request) {
	var rating models.RatingWisata
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

// DeleteRatingWisata untuk menghapus rating wisata berdasarkan ID
func DeleteRatingWisata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.RatingWisata{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
