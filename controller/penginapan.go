package controller

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"os"
	"time"
	"xplore/config"
	"xplore/models"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

// CreatePenginapan menerima request multi form-data
func CreatePenginapan(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	namaPenginapan := r.FormValue("nama_penginapan")
	lokasiPenginapan := r.FormValue("lokasi_penginapan")
	koordinatPenginapan := r.FormValue("koordinat_penginapan")
	kontakPenginapan := r.FormValue("kontak_penginapan")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["foto_penginapan[]"]
	var fotoPenginapanLinks []string
	for _, fileHeader := range files {
		// Buka file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload ke GCS dan dapatkan link-nya
		link, err := uploadPenginapanToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoPenginapanLinks = append(fotoPenginapanLinks, link)
	}

	penginapan := models.Penginapan{
		FotoPenginapan:      fotoPenginapanLinks,
		NamaPenginapan:      namaPenginapan,
		LokasiPenginapan:    lokasiPenginapan,
		KoordinatPenginapan: koordinatPenginapan,
		KontakPenginapan:    kontakPenginapan,
		Deskripsi:           deskripsi,
		Fasilitas:           fasilitas, // Array fasilitas
	}

	if err := config.DB.Create(&penginapan).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(penginapan)
}

// uploadPenginapanToGCS mengupload file ke Google Cloud Storage untuk Penginapan
func uploadPenginapanToGCS(file multipart.File, filename string) (string, error) {
	var opt option.ClientOption
	if os.Getenv("NODE_ENV") == "LOCAL" {
		base64Json := os.Getenv("GCP_SERVICE_ACCOUNT")
		if base64Json == "" {
			return "", fmt.Errorf("GCP_SERVICE_ACCOUNT environment variable is not set")
		}
		
		decodedData, err := base64.StdEncoding.DecodeString(base64Json)
		if err != nil {
			return "", fmt.Errorf("error decoding base64: %v", err.Error())
		}
		opt = option.WithCredentialsJSON(decodedData)
	} else {
		opt = option.WithCredentials(nil)
	}
	ctx := context.Background()
	client, err := storage.NewClient(ctx, opt)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Menentukan bucket dan objek
	bucket := client.Bucket("xplores")
	object := bucket.Object(filename)

	// Membuka writer untuk meng-upload file
	writer := object.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to copy file to GCS: %v", err)
	}

	// Menutup writer dan menyelesaikan upload
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %v", err)
	}

	// Mengembalikan URL file yang telah di-upload
	return fmt.Sprintf("https://storage.googleapis.com/xplores/%s", filename), nil
}

// GetAllPenginapan untuk mendapatkan semua penginapan
func GetAllPenginapan(w http.ResponseWriter, r *http.Request) {
	var penginapans []models.Penginapan
	if result := config.DB.Find(&penginapans); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(penginapans)
}

// GetPenginapanByID untuk mendapatkan penginapan berdasarkan ID
func GetPenginapanByID(w http.ResponseWriter, r *http.Request) {
	var penginapan models.Penginapan
	id := mux.Vars(r)["id"]

	// Menyaring penginapan berdasarkan ID
	if result := config.DB.First(&penginapan, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(penginapan)
}

// UpdatePenginapan untuk memperbarui penginapan berdasarkan ID
func UpdatePenginapan(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil ID Penginapan dari URL params
	vars := mux.Vars(r)
	idPenginapan := vars["id"]

	// Ambil data dari form
	namaPenginapan := r.FormValue("nama_penginapan")
	lokasiPenginapan := r.FormValue("lokasi_penginapan")
	koordinatPenginapan := r.FormValue("koordinat_penginapan")
	kontakPenginapan := r.FormValue("kontak_penginapan")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	// Ambil file gambar baru
	files := r.MultipartForm.File["foto_penginapan[]"]
	var fotoPenginapanLinks []string

	// Jika ada foto baru, upload dan simpan link
	for _, fileHeader := range files {
		// Buka file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload ke GCS dan dapatkan link-nya
		link, err := uploadPenginapanToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoPenginapanLinks = append(fotoPenginapanLinks, link)
	}

	// Ambil data penginapan yang ada di database
	var penginapan models.Penginapan
	if err := config.DB.First(&penginapan, idPenginapan).Error; err != nil {
		http.Error(w, "Penginapan not found", http.StatusNotFound)
		return
	}

	// Update data penginapan
	penginapan.NamaPenginapan = namaPenginapan
	penginapan.LokasiPenginapan = lokasiPenginapan
	penginapan.KoordinatPenginapan = koordinatPenginapan
	penginapan.KontakPenginapan = kontakPenginapan
	penginapan.Deskripsi = deskripsi
	penginapan.Fasilitas = fasilitas

	// Tambahkan foto baru jika ada, namun foto lama tetap ada
	penginapan.FotoPenginapan = fotoPenginapanLinks // Gantikan foto lama dengan foto baru

	// Update timestamp
	penginapan.UpdatedAt = time.Now()

	// Simpan update ke database
	if err := config.DB.Save(&penginapan).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(penginapan)
}

// DeletePenginapan untuk menghapus penginapan berdasarkan ID
func DeletePenginapan(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.Penginapan{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchPenginapan untuk mencari penginapan berdasarkan nama
func SearchPenginapan(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var penginapans []models.Penginapan

	// Query pencarian berdasarkan 'nama_penginapan', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_penginapan ILIKE ?", "%"+searchQuery+"%").Find(&penginapans); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(penginapans)
}

// GetPenginapanWithPagination - Fungsi untuk menampilkan data penginapan dengan paginasi
func GetPenginapanWithPagination(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter 'page' dari query string
	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1" // Default ke halaman 1 jika tidak ada parameter 'page'
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	// Tentukan jumlah data per halaman (misalnya 10)
	perPage := 10

	// Hitung OFFSET berdasarkan halaman yang diminta
	offset := (page - 1) * perPage

	// Inisialisasi slice untuk menampung hasil penginapan
	var penginapans []models.Penginapan

	// Query untuk mengambil data penginapan dengan paginasi (LIMIT dan OFFSET)
	if result := config.DB.Offset(offset).Limit(perPage).Find(&penginapans); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(penginapans)
}

// Fungsi untuk mendapatkan rating penginapan dan jumlah pemberi rating
func GetAllPenginapanRating(w http.ResponseWriter, r *http.Request) {
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(penginapanWithRating)
}

// Fungsi untuk mencari penginapan berdasarkan nama, dan menghitung rating rata-rata serta jumlah pemberi rating
func SearchPenginapanRating(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var penginapans []models.Penginapan

	// Query pencarian berdasarkan 'nama_penginapan', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_penginapan ILIKE ?", "%"+searchQuery+"%").Find(&penginapans); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Menampung hasil pencarian dengan rating rata-rata dan jumlah pemberi rating
	var penginapanWithRating []map[string]interface{}
	for _, p := range penginapans {
		// Mengambil data rating untuk penginapan
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

		// Menambahkan data penginapan, rating rata-rata, dan jumlah pemberi rating
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(penginapanWithRating)
}

// Fungsi untuk mendapatkan detail penginapan berdasarkan id, dengan rata-rata rating, jumlah pemberi rating dan ulasan
func GetPenginapanDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data penginapan berdasarkan id
	var penginapan models.Penginapan
	if err := config.DB.Where("id_penginapan = ?", id).First(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan not found", http.StatusNotFound)
		return
	}

	// Kirimkan detail data penginapan terlebih dahulu
	response := map[string]interface{}{
		"penginapan": penginapan,
	}

	// Ambil data rating untuk penginapan
	var rating []models.RatingPenginapan
	if err := config.DB.Where("id_penginapan = ?", penginapan.IDPenginapan).Find(&rating).Error; err != nil {
		http.Error(w, "Data rating not found", http.StatusInternalServerError)
		return
	}

	// Menghitung rata-rata rating dan jumlah pemberi rating per nilai
	var averageRating float64
	var ratingCounts = map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	var totalRating float64

	for _, r := range rating {
		totalRating += r.Rating
		ratingCounts[int(r.Rating)]++ // Menghitung jumlah pemberi rating untuk setiap nilai
	}

	// Menghitung rata-rata rating
	if len(rating) > 0 {
		averageRating = totalRating / float64(len(rating))
	}

	// Menghitung jumlah pemberi rating unik
	var uniqueRaters int64
	config.DB.Model(&models.RatingPenginapan{}).Where("id_penginapan = ?", penginapan.IDPenginapan).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response["rating_average"] = averageRating
	response["rating_counts"] = ratingCounts
	response["reviewers_count"] = uniqueRaters
	response["ratings"] = rating

	// Menambahkan data penginapan lainnya
	var penginapanLain []models.Penginapan
	if err := config.DB.Where("id_penginapan != ?", id).Find(&penginapanLain).Error; err != nil {
		http.Error(w, "Data penginapan lain not found", http.StatusInternalServerError)
		return
	}

	var penginapanLainWithDetails []map[string]interface{}
	for _, p := range penginapanLain {
		var penginapanData = map[string]interface{}{
			"id_penginapan":         p.IDPenginapan,
			"nama_penginapan":       p.NamaPenginapan,
			"lokasi_penginapan":     p.LokasiPenginapan,
			"koordinat_penginapan":  p.KoordinatPenginapan,
			"kontak_penginapan":     p.KontakPenginapan,
			"foto_penginapan":       p.FotoPenginapan,
			"deskripsi":             p.Deskripsi,
			"fasilitas":             p.Fasilitas,
			"rating":                averageRating, // Set default average rating (for demo)
			"jumlah_pemberi_rating": uniqueRaters,  // Set default reviewer count
		}
		penginapanLainWithDetails = append(penginapanLainWithDetails, penginapanData)
	}

	// Menambahkan data penginapan lainnya ke dalam response
	response["penginapan_lainnnya"] = penginapanLainWithDetails

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi untuk mendapatkan detail rating penginapan berdasarkan id, dengan rata-rata rating, jumlah pemberi rating, dan ulasan
func GetPenginapanRatingDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data penginapan berdasarkan id
	var penginapan models.Penginapan
	if err := config.DB.Where("id_penginapan = ?", id).First(&penginapan).Error; err != nil {
		http.Error(w, "Data penginapan not found", http.StatusNotFound)
		return
	}

	// Kirimkan detail data penginapan terlebih dahulu
	response := map[string]interface{}{
		"penginapan": penginapan,
	}

	// Ambil data rating untuk penginapan
	var rating []models.RatingPenginapan
	if err := config.DB.Where("id_penginapan = ?", penginapan.IDPenginapan).Find(&rating).Error; err != nil {
		http.Error(w, "Data rating not found", http.StatusInternalServerError)
		return
	}

	// Menghitung rata-rata rating dan jumlah pemberi rating per nilai
	var averageRating float64
	var ratingCounts = map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0}
	var totalRating float64

	for _, r := range rating {
		totalRating += r.Rating
		ratingCounts[int(r.Rating)]++ // Menghitung jumlah pemberi rating untuk setiap nilai
	}

	// Menghitung rata-rata rating
	if len(rating) > 0 {
		averageRating = totalRating / float64(len(rating))
	}

	// Menghitung jumlah pemberi rating unik
	var uniqueRaters int64
	config.DB.Model(&models.RatingPenginapan{}).Where("id_penginapan = ?", penginapan.IDPenginapan).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response["rating_average"] = averageRating
	response["rating_counts"] = ratingCounts
	response["reviewers_count"] = uniqueRaters
	response["ratings"] = rating

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
