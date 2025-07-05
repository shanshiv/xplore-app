package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"
)

// Fungsi untuk mendapatkan dashboard yang berisi kuliner, penginapan, wisata, dan jumlah data
func GetDashboard(w http.ResponseWriter, r *http.Request) {
	// Ambil data Kuliner
	var kuliner []models.Kuliner
	if err := config.DB.Find(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner not found", http.StatusNotFound)
		return
	}

	var kulinerWithRating []map[string]interface{}
	for _, k := range kuliner {
		var rating []models.RatingKuliner
		if err := config.DB.Where("id_kuliner = ?", k.IDKuliner).Find(&rating).Error; err != nil {
			rating = []models.RatingKuliner{} // Jika tidak ada rating
		}

		var averageRating float64
		if len(rating) > 0 {
			var totalRating float64
			for _, r := range rating {
				totalRating += r.Rating
			}
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		kulinerData := map[string]interface{}{
			"id_kuliner":        k.IDKuliner,
			"nama_kuliner":      k.NamaKuliner,
			"lokasi_kuliner":    k.LokasiKuliner,
			"koordinat_kuliner": k.KoordinatKuliner,
			"kontak_kuliner":    k.KontakKuliner,
			"foto_kuliner":      k.FotoKuliner,
			"deskripsi":         k.Deskripsi,
			"fasilitas":         k.Fasilitas,
			"rating":            averageRating,
		}

		kulinerWithRating = append(kulinerWithRating, kulinerData)
	}

	// Ambil data Penginapan
	var penginapan []models.Penginapan
	if err := config.DB.Find(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan not found", http.StatusNotFound)
		return
	}

	var penginapanWithRating []map[string]interface{}
	for _, p := range penginapan {
		var rating []models.RatingPenginapan
		if err := config.DB.Where("id_penginapan = ?", p.IDPenginapan).Find(&rating).Error; err != nil {
			rating = []models.RatingPenginapan{} // Jika tidak ada rating
		}

		var averageRating float64
		if len(rating) > 0 {
			var totalRating float64
			for _, r := range rating {
				totalRating += r.Rating
			}
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		penginapanData := map[string]interface{}{
			"id_penginapan":        p.IDPenginapan,
			"nama_penginapan":      p.NamaPenginapan,
			"lokasi_penginapan":    p.LokasiPenginapan,
			"koordinat_penginapan": p.KoordinatPenginapan,
			"kontak_penginapan":    p.KontakPenginapan,
			"foto_penginapan":      p.FotoPenginapan,
			"deskripsi":            p.Deskripsi,
			"fasilitas":            p.Fasilitas,
			"rating":               averageRating,
		}

		penginapanWithRating = append(penginapanWithRating, penginapanData)
	}

	// Ambil data Wisata
	var wisata []models.Wisata
	if err := config.DB.Find(&wisata).Error; err != nil {
		http.Error(w, "Data wisata not found", http.StatusNotFound)
		return
	}

	var wisataWithRating []map[string]interface{}
	for _, w := range wisata {
		var rating []models.RatingWisata
		if err := config.DB.Where("id_wisata = ?", w.IDWisata).Find(&rating).Error; err != nil {
			rating = []models.RatingWisata{} // Jika tidak ada rating
		}

		var averageRating float64
		if len(rating) > 0 {
			var totalRating float64
			for _, r := range rating {
				totalRating += r.Rating
			}
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		wisataData := map[string]interface{}{
			"id_wisata":        w.IDWisata,
			"nama_wisata":      w.NamaWisata,
			"lokasi_wisata":    w.LokasiWisata,
			"koordinat_wisata": w.KoordinatWisata,
			"kontak_wisata":    w.KontakWisata,
			"foto_wisata":      w.FotoWisata,
			"deskripsi":        w.Deskripsi,
			"fasilitas":        w.Fasilitas,
			"rating":           averageRating,
		}

		wisataWithRating = append(wisataWithRating, wisataData)
	}

	// Menghitung jumlah data
	var kulinerCount int64
	var penginapanCount int64
	var wisataCount int64

	// Hitung jumlah data untuk masing-masing tabel
	config.DB.Model(&models.Kuliner{}).Count(&kulinerCount)
	config.DB.Model(&models.Penginapan{}).Count(&penginapanCount)
	config.DB.Model(&models.Wisata{}).Count(&wisataCount)

	// Gabungkan data Kuliner, Penginapan, Wisata dan jumlah data dalam satu response
	response := map[string]interface{}{
		"kuliner":          kulinerWithRating,
		"penginapan":       penginapanWithRating,
		"wisata":           wisataWithRating,
		"kuliner_count":    kulinerCount,
		"penginapan_count": penginapanCount,
		"wisata_count":     wisataCount,
	}

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
