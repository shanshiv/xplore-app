package controller

import (
	"encoding/json"
	"net/http"
	"xplore/config"
	"xplore/models"
)

// Fungsi untuk mendapatkan data beranda (wisata, kuliner, penginapan) dan rata-rata rating serta jumlah pemberi rating
func GetBeranda(w http.ResponseWriter, r *http.Request) {
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
		var totalRating float64
		for _, r := range rating {
			totalRating += r.Rating
		}
		if len(rating) > 0 {
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		// Menghitung jumlah pemberi rating unik
		var uniqueRaters int64
		config.DB.Model(&models.RatingWisata{}).Where("id_wisata = ?", w.IDWisata).Distinct("id_akun").Count(&uniqueRaters)

		wisataData := map[string]interface{}{
			"id_wisata":             w.IDWisata,
			"nama_wisata":           w.NamaWisata,
			"lokasi_wisata":         w.LokasiWisata,
			"koordinat_wisata":      w.KoordinatWisata,
			"kontak_wisata":         w.KontakWisata,
			"foto_wisata":           w.FotoWisata,
			"deskripsi":             w.Deskripsi,
			"fasilitas":             w.Fasilitas,
			"rating":                averageRating,
			"jumlah_pemberi_rating": uniqueRaters, // Menambahkan jumlah pemberi rating
		}

		wisataWithRating = append(wisataWithRating, wisataData)
	}

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
		var totalRating float64
		for _, r := range rating {
			totalRating += r.Rating
		}
		if len(rating) > 0 {
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		// Menghitung jumlah pemberi rating unik
		var uniqueRaters int64
		config.DB.Model(&models.RatingKuliner{}).Where("id_kuliner = ?", k.IDKuliner).Distinct("id_akun").Count(&uniqueRaters)

		kulinerData := map[string]interface{}{
			"id_kuliner":            k.IDKuliner,
			"nama_kuliner":          k.NamaKuliner,
			"lokasi_kuliner":        k.LokasiKuliner,
			"koordinat_kuliner":     k.KoordinatKuliner,
			"kontak_kuliner":        k.KontakKuliner,
			"foto_kuliner":          k.FotoKuliner,
			"deskripsi":             k.Deskripsi,
			"fasilitas":             k.Fasilitas,
			"rating":                averageRating,
			"jumlah_pemberi_rating": uniqueRaters, // Menambahkan jumlah pemberi rating
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
		var totalRating float64
		for _, r := range rating {
			totalRating += r.Rating
		}
		if len(rating) > 0 {
			averageRating = totalRating / float64(len(rating))
		} else {
			averageRating = 0
		}

		// Menghitung jumlah pemberi rating unik
		var uniqueRaters int64
		config.DB.Model(&models.RatingPenginapan{}).Where("id_penginapan = ?", p.IDPenginapan).Distinct("id_akun").Count(&uniqueRaters)

		penginapanData := map[string]interface{}{
			"id_penginapan":         p.IDPenginapan,
			"nama_penginapan":       p.NamaPenginapan,
			"lokasi_penginapan":     p.LokasiPenginapan,
			"koordinat_penginapan":  p.KoordinatPenginapan,
			"kontak_penginapan":     p.KontakPenginapan,
			"foto_penginapan":       p.FotoPenginapan,
			"deskripsi":             p.Deskripsi,
			"fasilitas":             p.Fasilitas,
			"rating":                averageRating,
			"jumlah_pemberi_rating": uniqueRaters, // Menambahkan jumlah pemberi rating
		}

		penginapanWithRating = append(penginapanWithRating, penginapanData)
	}

	// Gabungkan data Wisata, Kuliner, Penginapan dalam satu response
	response := map[string]interface{}{
		"wisata":     wisataWithRating,
		"kuliner":    kulinerWithRating,
		"penginapan": penginapanWithRating,
	}

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
