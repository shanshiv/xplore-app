package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"
)

// Fungsi untuk mengambil semua data kuliner, penginapan, dan wisata
func GetPeta(w http.ResponseWriter, r *http.Request) {
	// Ambil semua data kuliner
	var kuliner []models.Kuliner
	if err := config.DB.Find(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner not found", http.StatusNotFound)
		return
	}

	// Ambil semua data penginapan
	var penginapan []models.Penginapan
	if err := config.DB.Find(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan not found", http.StatusNotFound)
		return
	}

	// Ambil semua data wisata
	var wisata []models.Wisata
	if err := config.DB.Find(&wisata).Error; err != nil {
		http.Error(w, "Data wisata not found", http.StatusNotFound)
		return
	}

	// Menyiapkan response data
	response := map[string]interface{}{
		"kuliner":    kuliner,
		"penginapan": penginapan,
		"wisata":     wisata,
	}

	// Kirimkan response ke client dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
