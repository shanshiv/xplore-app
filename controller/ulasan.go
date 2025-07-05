package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"
)

// Fungsi untuk mengambil semua ulasan dari kuliner, penginapan, dan wisata
func GetAllRatings(w http.ResponseWriter, r *http.Request) {
	// Ambil ulasan kuliner
	var kulinerRatings []models.RatingKuliner
	if err := config.DB.Find(&kulinerRatings).Error; err != nil {
		http.Error(w, "Data rating kuliner not found", http.StatusNotFound)
		return
	}

	// Ambil ulasan penginapan
	var penginapanRatings []models.RatingPenginapan
	if err := config.DB.Find(&penginapanRatings).Error; err != nil {
		http.Error(w, "Data rating penginapan not found", http.StatusNotFound)
		return
	}

	// Ambil ulasan wisata
	var wisataRatings []models.RatingWisata
	if err := config.DB.Find(&wisataRatings).Error; err != nil {
		http.Error(w, "Data rating wisata not found", http.StatusNotFound)
		return
	}

	// Gabungkan ketiga data ulasan
	response := map[string]interface{}{
		"kuliner":    kulinerRatings,
		"penginapan": penginapanRatings,
		"wisata":     wisataRatings,
	}

	// Kirimkan data ulasan ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
