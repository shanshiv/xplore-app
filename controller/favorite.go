package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"

	"github.com/gorilla/mux"
)

// Fungsi untuk mendapatkan semua favorit kuliner, penginapan, dan wisata berdasarkan user_id
func GetAllFavorit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"] // Mendapatkan user_id dari URL parameter

	// Ambil semua kuliner favorit milik user
	var userFavoritesKuliner []models.UserFavoritesKuliner
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavoritesKuliner).Error; err != nil {
		http.Error(w, "Data kuliner favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua kuliner berdasarkan kuliner_id yang ada pada userFavorites
	var kuliner []models.Kuliner
	var kulinerIds []uint
	for _, uf := range userFavoritesKuliner {
		kulinerIds = append(kulinerIds, uf.KulinerId) // Menyusun daftar kuliner_id
	}

	if err := config.DB.Where("id_kuliner IN (?)", kulinerIds).Find(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk kuliner
	var kulinerWithAvgRating []map[string]interface{}
	for _, k := range kuliner {
		var ratings []models.RatingKuliner
		if err := config.DB.Where("id_kuliner = ?", k.IDKuliner).Find(&ratings).Error; err != nil {
			ratings = []models.RatingKuliner{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}
		kulinerData := map[string]interface{}{
			"kuliner":        k,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		kulinerWithAvgRating = append(kulinerWithAvgRating, kulinerData)
	}

	// Ambil data favorit penginapan milik user
	var userFavoritesPenginapan []models.UserFavoritesPenginapan
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavoritesPenginapan).Error; err != nil {
		http.Error(w, "Data penginapan favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua penginapan berdasarkan penginapan_id yang ada pada userFavorites
	var penginapan []models.Penginapan
	var penginapanIds []uint
	for _, uf := range userFavoritesPenginapan {
		penginapanIds = append(penginapanIds, uf.PenginapanId) // Menyusun daftar penginapan_id
	}

	if err := config.DB.Where("id_penginapan IN (?)", penginapanIds).Find(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk penginapan
	var penginapanWithAvgRating []map[string]interface{}
	for _, p := range penginapan {
		var ratings []models.RatingPenginapan
		if err := config.DB.Where("id_penginapan = ?", p.IDPenginapan).Find(&ratings).Error; err != nil {
			ratings = []models.RatingPenginapan{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}
		penginapanData := map[string]interface{}{
			"penginapan":     p,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		penginapanWithAvgRating = append(penginapanWithAvgRating, penginapanData)
	}

	// Ambil data favorit wisata milik user
	var userFavoritesWisata []models.UserFavoritesWisata
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavoritesWisata).Error; err != nil {
		http.Error(w, "Data wisata favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua wisata berdasarkan wisata_id yang ada pada userFavorites
	var wisata []models.Wisata
	var wisataIds []uint
	for _, uf := range userFavoritesWisata {
		wisataIds = append(wisataIds, uf.WisataId) // Menyusun daftar wisata_id
	}

	if err := config.DB.Where("id_wisata IN (?)", wisataIds).Find(&wisata).Error; err != nil {
		http.Error(w, "Data wisata tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk wisata
	var wisataWithAvgRating []map[string]interface{}
	for _, w := range wisata {
		var ratings []models.RatingWisata
		if err := config.DB.Where("id_wisata = ?", w.IDWisata).Find(&ratings).Error; err != nil {
			ratings = []models.RatingWisata{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}
		wisataData := map[string]interface{}{
			"wisata":         w,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		wisataWithAvgRating = append(wisataWithAvgRating, wisataData)
	}

	// Menyiapkan response data
	response := map[string]interface{}{
		"kuliners":    kulinerWithAvgRating,
		"penginapans": penginapanWithAvgRating,
		"wisatas":     wisataWithAvgRating,
	}

	// Kirimkan response ke client dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi untuk menambahkan favorit kuliner bagi user
func AddKulinerFavorite(w http.ResponseWriter, r *http.Request) {
	var favorite models.UserFavoritesKuliner
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&favorite); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", favorite.UserID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menyimpan data favorit kuliner
	if result := config.DB.Create(&favorite); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kuliner favorite added successfully"})
}

// Fungsi untuk menambahkan favorit penginapan bagi user
func AddPenginapanFavorite(w http.ResponseWriter, r *http.Request) {
	var favorite models.UserFavoritesPenginapan
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&favorite); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", favorite.UserID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menyimpan data favorit penginapan
	if result := config.DB.Create(&favorite); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Penginapan favorite added successfully"})
}

// Fungsi untuk menambahkan favorit wisata bagi user
func AddWisataFavorite(w http.ResponseWriter, r *http.Request) {
	var favorite models.UserFavoritesWisata
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&favorite); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", favorite.UserID).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menyimpan data favorit wisata
	if result := config.DB.Create(&favorite); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Wisata favorite added successfully"})
}

// Fungsi untuk menghapus favorit kuliner bagi user
func DeleteKulinerFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]
	kulinerId := vars["kuliner_id"]

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menghapus data favorit kuliner
	if result := config.DB.Where("user_id = ? AND kuliner_id = ?", userId, kulinerId).Delete(&models.UserFavoritesKuliner{}); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kuliner favorite deleted successfully"})
}

// Fungsi untuk menghapus favorit penginapan bagi user
func DeletePenginapanFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]
	penginapanId := vars["penginapan_id"]

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menghapus data favorit penginapan
	if result := config.DB.Where("user_id = ? AND penginapan_id = ?", userId, penginapanId).Delete(&models.UserFavoritesPenginapan{}); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Penginapan favorite deleted successfully"})
}

// Fungsi untuk menghapus favorit wisata bagi user
func DeleteWisataFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["user_id"]
	wisataId := vars["wisata_id"]

	// Cek apakah user ada berdasarkan user_id
	var user models.Akun
	if err := config.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Menghapus data favorit wisata
	if result := config.DB.Where("user_id = ? AND wisata_id = ?", userId, wisataId).Delete(&models.UserFavoritesWisata{}); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Wisata favorite deleted successfully"})
}

// Fungsi untuk mendapatkan semua favorit kuliner berdasarkan user_id
func GetAllKulinerFavorit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"] // Mendapatkan user_id dari URL parameter

	// Ambil semua data favorit kuliner berdasarkan user_id
	var userFavorites []models.UserFavoritesKuliner
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavorites).Error; err != nil {
		http.Error(w, "Data kuliner favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua kuliner berdasarkan kuliner_id yang ada pada userFavorites
	var kuliner []models.Kuliner
	var kulinerIds []uint
	for _, uf := range userFavorites {
		kulinerIds = append(kulinerIds, uf.KulinerId) // Menyusun daftar kuliner_id
	}

	if err := config.DB.Where("id_kuliner IN (?)", kulinerIds).Find(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk kuliner
	var kulinerWithAvgRating []map[string]interface{}
	for _, k := range kuliner {
		var ratings []models.RatingKuliner
		if err := config.DB.Where("id_kuliner = ?", k.IDKuliner).Find(&ratings).Error; err != nil {
			ratings = []models.RatingKuliner{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}

		kulinerData := map[string]interface{}{
			"kuliner":        k,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		kulinerWithAvgRating = append(kulinerWithAvgRating, kulinerData)
	}

	// Kirimkan response ke client dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kulinerWithAvgRating)
}

// Fungsi untuk mendapatkan semua favorit penginapan berdasarkan user_id
func GetAllPenginapanFavorit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"] // Mendapatkan user_id dari URL parameter

	// Ambil semua data favorit penginapan berdasarkan user_id
	var userFavorites []models.UserFavoritesPenginapan
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavorites).Error; err != nil {
		http.Error(w, "Data penginapan favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua penginapan berdasarkan penginapan_id yang ada pada userFavorites
	var penginapan []models.Penginapan
	var penginapanIds []uint
	for _, uf := range userFavorites {
		penginapanIds = append(penginapanIds, uf.PenginapanId) // Menyusun daftar penginapan_id
	}

	if err := config.DB.Where("id_penginapan IN (?)", penginapanIds).Find(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk penginapan
	var penginapanWithAvgRating []map[string]interface{}
	for _, p := range penginapan {
		var ratings []models.RatingPenginapan
		if err := config.DB.Where("id_penginapan = ?", p.IDPenginapan).Find(&ratings).Error; err != nil {
			ratings = []models.RatingPenginapan{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}

		penginapanData := map[string]interface{}{
			"penginapan":     p,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		penginapanWithAvgRating = append(penginapanWithAvgRating, penginapanData)
	}

	// Kirimkan response ke client dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(penginapanWithAvgRating)
}

// Fungsi untuk mendapatkan semua favorit wisata berdasarkan user_id
func GetAllWisataFavorit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["id"] // Mendapatkan user_id dari URL parameter

	// Ambil semua data favorit wisata berdasarkan user_id
	var userFavorites []models.UserFavoritesWisata
	if err := config.DB.Where("user_id = ?", userId).Find(&userFavorites).Error; err != nil {
		http.Error(w, "Data wisata favorit tidak ditemukan", http.StatusNotFound)
		return
	}

	// Ambil semua wisata berdasarkan wisata_id yang ada pada userFavorites
	var wisata []models.Wisata
	var wisataIds []uint
	for _, uf := range userFavorites {
		wisataIds = append(wisataIds, uf.WisataId) // Menyusun daftar wisata_id
	}

	if err := config.DB.Where("id_wisata IN (?)", wisataIds).Find(&wisata).Error; err != nil {
		http.Error(w, "Data wisata tidak ditemukan", http.StatusNotFound)
		return
	}

	// Menghitung rata-rata rating untuk wisata
	var wisataWithAvgRating []map[string]interface{}
	for _, w := range wisata {
		var ratings []models.RatingWisata
		if err := config.DB.Where("id_wisata = ?", w.IDWisata).Find(&ratings).Error; err != nil {
			ratings = []models.RatingWisata{} // Jika tidak ada rating
		}

		var totalRating float64
		for _, r := range ratings {
			totalRating += r.Rating
		}

		var averageRating float64
		if len(ratings) > 0 {
			averageRating = totalRating / float64(len(ratings))
		}

		wisataData := map[string]interface{}{
			"wisata":         w,
			"rating_average": averageRating,
			"rating_count":   len(ratings),
		}

		wisataWithAvgRating = append(wisataWithAvgRating, wisataData)
	}

	// Kirimkan response ke client dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wisataWithAvgRating)
}
