package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
	"xplore/config"
	"xplore/models"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

// CreateKuliner menerima request multi form-data
func CreateKuliner(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	namaKuliner := r.FormValue("nama_kuliner")
	lokasiKuliner := r.FormValue("lokasi_kuliner")
	koordinatKuliner := r.FormValue("koordinat_kuliner")
	kontakKuliner := r.FormValue("kontak_kuliner")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["foto_kuliner[]"]
	var fotoKulinerLinks []string
	for _, fileHeader := range files {
		// Buka file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload ke GCS dan dapatkan link-nya
		link, err := uploadKulinerToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoKulinerLinks = append(fotoKulinerLinks, link)
	}

	kuliner := models.Kuliner{
		FotoKuliner:      fotoKulinerLinks,
		NamaKuliner:      namaKuliner,
		LokasiKuliner:    lokasiKuliner,
		KoordinatKuliner: koordinatKuliner,
		KontakKuliner:    kontakKuliner,
		Deskripsi:        deskripsi,
		Fasilitas:        fasilitas, // Array fasilitas
	}

	if err := config.DB.Create(&kuliner).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(kuliner)
}

// uploadKulinerToGCS mengupload file ke Google Cloud Storage untuk Kuliner
func uploadKulinerToGCS(file multipart.File, filename string) (string, error) {
	ctx := context.Background()

	// Membuat klien Google Cloud Storage
	client, err := storage.NewClient(ctx, option.WithCredentialsFile("./config/xplore-48-447519269b91.json"))
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

// GetAllKuliner untuk mendapatkan semua kuliner
func GetAllKuliner(w http.ResponseWriter, r *http.Request) {
	var kuliners []models.Kuliner
	if result := config.DB.Find(&kuliners); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kuliners)
}

// GetKulinerByID untuk mendapatkan kuliner berdasarkan ID
func GetKulinerByID(w http.ResponseWriter, r *http.Request) {
	var kuliner models.Kuliner
	id := mux.Vars(r)["id"]

	// Menyaring kuliner berdasarkan ID
	if result := config.DB.First(&kuliner, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kuliner)
}

// UpdateKuliner untuk memperbarui kuliner berdasarkan ID
func UpdateKuliner(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil ID Kuliner dari URL params
	vars := mux.Vars(r)
	idKuliner := vars["id"]

	// Ambil data dari form
	namaKuliner := r.FormValue("nama_kuliner")
	lokasiKuliner := r.FormValue("lokasi_kuliner")
	koordinatKuliner := r.FormValue("koordinat_kuliner")
	kontakKuliner := r.FormValue("kontak_kuliner")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	// Ambil file gambar baru
	files := r.MultipartForm.File["foto_kuliner[]"]
	var fotoKulinerLinks []string

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
		link, err := uploadKulinerToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoKulinerLinks = append(fotoKulinerLinks, link)
	}

	// Ambil data kuliner yang ada di database
	var kuliner models.Kuliner
	if err := config.DB.First(&kuliner, idKuliner).Error; err != nil {
		http.Error(w, "Kuliner not found", http.StatusNotFound)
		return
	}

	// Update data kuliner
	kuliner.NamaKuliner = namaKuliner
	kuliner.LokasiKuliner = lokasiKuliner
	kuliner.KoordinatKuliner = koordinatKuliner
	kuliner.KontakKuliner = kontakKuliner
	kuliner.Deskripsi = deskripsi
	kuliner.Fasilitas = fasilitas

	// Tambahkan foto baru jika ada, namun foto lama tetap ada
	kuliner.FotoKuliner = fotoKulinerLinks // Gantikan foto lama dengan foto baru

	// Update timestamp
	kuliner.UpdatedAt = time.Now()

	// Simpan update ke database
	if err := config.DB.Save(&kuliner).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kuliner)
}

// DeleteKuliner untuk menghapus kuliner berdasarkan ID
func DeleteKuliner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.Kuliner{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchKuliner untuk mencari kuliner berdasarkan nama
func SearchKuliner(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var kuliners []models.Kuliner

	// Query pencarian berdasarkan 'nama_kuliner', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_kuliner ILIKE ?", "%"+searchQuery+"%").Find(&kuliners); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kuliners)
}

// GetKulinerWithPagination - Fungsi untuk menampilkan data kuliner dengan paginasi
func GetKulinerWithPagination(w http.ResponseWriter, r *http.Request) {
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

	// Inisialisasi slice untuk menampung hasil kuliner
	var kuliners []models.Kuliner

	// Query untuk mengambil data kuliner dengan paginasi (LIMIT dan OFFSET)
	if result := config.DB.Offset(offset).Limit(perPage).Find(&kuliners); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kuliners)
}

// Fungsi untuk mendapatkan rating kuliner dan jumlah pemberi rating
func GetAllKulinerRating(w http.ResponseWriter, r *http.Request) {
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kulinerWithRating)
}

// Fungsi untuk mencari kuliner berdasarkan nama, dan menghitung rating rata-rata serta jumlah pemberi rating
func SearchKulinerRating(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var kuliners []models.Kuliner

	// Query pencarian berdasarkan 'nama_kuliner', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_kuliner ILIKE ?", "%"+searchQuery+"%").Find(&kuliners); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Menampung hasil pencarian dengan rating rata-rata dan jumlah pemberi rating
	var kulinerWithRating []map[string]interface{}
	for _, k := range kuliners {
		// Mengambil data rating untuk kuliner
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

		// Menambahkan data kuliner, rating rata-rata, dan jumlah pemberi rating
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kulinerWithRating)
}

// Fungsi untuk mendapatkan detail kuliner berdasarkan id, dengan rata-rata rating, jumlah pemberi rating dan ulasan
func GetKulinerDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data kuliner berdasarkan id
	var kuliner models.Kuliner
	if err := config.DB.Where("id_kuliner = ?", id).First(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner not found", http.StatusNotFound)
		return
	}

	// Ambil data rating untuk kuliner
	var rating []models.RatingKuliner
	if err := config.DB.Where("id_kuliner = ?", kuliner.IDKuliner).Find(&rating).Error; err != nil {
		http.Error(w, "Data rating not found", http.StatusInternalServerError)
		return
	}

	// Inisialisasi variabel untuk menghitung rata-rata rating dan jumlah pemberi rating per nilai
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
	config.DB.Model(&models.RatingKuliner{}).Where("id_kuliner = ?", kuliner.IDKuliner).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response := map[string]interface{}{
		"kuliner":         kuliner,
		"rating_average":  averageRating,
		"rating_counts":   ratingCounts,
		"reviewers_count": uniqueRaters,
		"ratings":         rating,
	}

	// Menambahkan data kuliner lainnya (contoh tambahan kuliner)
	var kulinerLain []models.Kuliner
	if err := config.DB.Where("id_kuliner != ?", id).Find(&kulinerLain).Error; err != nil {
		http.Error(w, "Data kuliner lain not found", http.StatusInternalServerError)
		return
	}

	var kulinerLainWithDetails []map[string]interface{}
	for _, k := range kulinerLain {
		var kulinerData = map[string]interface{}{
			"id_kuliner":            k.IDKuliner,
			"nama_kuliner":          k.NamaKuliner,
			"lokasi_kuliner":        k.LokasiKuliner,
			"koordinat_kuliner":     k.KoordinatKuliner,
			"kontak_kuliner":        k.KontakKuliner,
			"foto_kuliner":          k.FotoKuliner,
			"deskripsi":             k.Deskripsi,
			"fasilitas":             k.Fasilitas,
			"rating":                averageRating, // Set default average rating (for demo)
			"jumlah_pemberi_rating": uniqueRaters,  // Set default reviewer count
		}
		kulinerLainWithDetails = append(kulinerLainWithDetails, kulinerData)
	}

	// Menambahkan data kuliner lainnya ke dalam response
	response["kuliner_lainnnya"] = kulinerLainWithDetails

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi untuk mendapatkan detail rating kuliner berdasarkan id, dengan rata-rata rating, jumlah pemberi rating, dan ulasan
func GetKulinerRatingDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data kuliner berdasarkan id
	var kuliner models.Kuliner
	if err := config.DB.Where("id_kuliner = ?", id).First(&kuliner).Error; err != nil {
		http.Error(w, "Data kuliner not found", http.StatusNotFound)
		return
	}

	// Kirimkan detail data kuliner terlebih dahulu
	response := map[string]interface{}{
		"kuliner": kuliner,
	}

	// Ambil data rating untuk kuliner
	var rating []models.RatingKuliner
	if err := config.DB.Where("id_kuliner = ?", kuliner.IDKuliner).Find(&rating).Error; err != nil {
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
	config.DB.Model(&models.RatingKuliner{}).Where("id_kuliner = ?", kuliner.IDKuliner).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response["rating_average"] = averageRating
	response["rating_counts"] = ratingCounts
	response["reviewers_count"] = uniqueRaters
	response["ratings"] = rating

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
