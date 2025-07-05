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

// CreateWisata menerima request multi form-data
func CreateWisata(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	namaWisata := r.FormValue("nama_wisata")
	lokasiWisata := r.FormValue("lokasi_wisata")
	koordinatWisata := r.FormValue("koordinat_wisata")
	kontakWisata := r.FormValue("kontak_wisata")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["foto_wisata[]"]
	var fotoWisataLinks []string
	for _, fileHeader := range files {
		// Buka file
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Unable to open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Upload ke GCS dan dapatkan link-nya
		link, err := uploadToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoWisataLinks = append(fotoWisataLinks, link)
	}
	wisata := models.Wisata{
		FotoWisata:      fotoWisataLinks,
		NamaWisata:      namaWisata,
		LokasiWisata:    lokasiWisata,
		KoordinatWisata: koordinatWisata, // Pastikan data koordinat diproses dengan benar
		KontakWisata:    kontakWisata,
		Deskripsi:       deskripsi,
		Fasilitas:       fasilitas, // Array fasilitas
		// Kontak wisata juga harus diproses dengan benar
		// Array URL foto
	}

	if err := config.DB.Create(&wisata).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(wisata)
}

// uploadToGCS mengupload file ke Google Cloud Storage
func uploadToGCS(file multipart.File, filename string) (string, error) {
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

// GetAllWisata untuk mendapatkan semua wisata
func GetAllWisata(w http.ResponseWriter, r *http.Request) {
	var wisatas []models.Wisata
	if result := config.DB.Find(&wisatas); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wisatas)
}

// GetWisataByID untuk mendapatkan wisata berdasarkan ID
func GetWisataByID(w http.ResponseWriter, r *http.Request) {
	var wisata models.Wisata
	id := mux.Vars(r)["id"]

	// Menyaring wisata berdasarkan ID
	if result := config.DB.First(&wisata, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wisata)
}

// UpdateWisata untuk memperbarui wisata berdasarkan ID
func UpdateWisata(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil ID Wisata dari URL params
	vars := mux.Vars(r)
	idWisata := vars["id"]

	// Ambil data dari form
	namaWisata := r.FormValue("nama_wisata")
	lokasiWisata := r.FormValue("lokasi_wisata")
	koordinatWisata := r.FormValue("koordinat_wisata")
	kontakWisata := r.FormValue("kontak_wisata")
	deskripsi := r.FormValue("deskripsi")
	fasilitasData := r.FormValue("fasilitas[]")
	var fasilitas []string
	if err := json.Unmarshal([]byte(fasilitasData), &fasilitas); err != nil {
		http.Error(w, "Unable to parse fasilitas", http.StatusBadRequest)
		return
	}

	// Ambil file gambar baru
	files := r.MultipartForm.File["foto_wisata[]"]
	var fotoWisataLinks []string

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
		link, err := uploadToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoWisataLinks = append(fotoWisataLinks, link)
	}

	// Ambil data wisata yang ada di database
	var wisata models.Wisata
	if err := config.DB.First(&wisata, idWisata).Error; err != nil {
		http.Error(w, "Wisata not found", http.StatusNotFound)
		return
	}

	// Update data wisata
	wisata.NamaWisata = namaWisata
	wisata.LokasiWisata = lokasiWisata
	wisata.KoordinatWisata = koordinatWisata
	wisata.KontakWisata = kontakWisata
	wisata.Deskripsi = deskripsi
	wisata.Fasilitas = fasilitas

	// Tambahkan foto baru jika ada, namun foto lama tetap ada
	wisata.FotoWisata = fotoWisataLinks // Gantikan foto lama dengan foto baru

	// Update timestamp
	wisata.UpdatedAt = time.Now()

	// Simpan update ke database
	if err := config.DB.Save(&wisata).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wisata)
}

// DeleteWisata untuk menghapus wisata berdasarkan ID
func DeleteWisata(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.Wisata{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchWisata untuk mencari wisata berdasarkan nama
func SearchWisata(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var wisatas []models.Wisata

	// Query pencarian berdasarkan 'nama_wisata', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_wisata ILIKE ?", "%"+searchQuery+"%").Find(&wisatas); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wisatas)
}

// GetWisataWithPagination - Fungsi untuk menampilkan data wisata dengan paginasi
func GetWisataWithPagination(w http.ResponseWriter, r *http.Request) {
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

	// Inisialisasi slice untuk menampung hasil wisata
	var wisatas []models.Wisata

	// Query untuk mengambil data wisata dengan paginasi (LIMIT dan OFFSET)
	if result := config.DB.Offset(offset).Limit(perPage).Find(&wisatas); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Kirimkan hasil pencarian dalam bentuk JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wisatas)
}

func GetAllWisataRating(w http.ResponseWriter, r *http.Request) {
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wisataWithRating)
}

// Fungsi untuk mencari wisata berdasarkan nama, dan menghitung rating rata-rata serta jumlah pemberi rating
func SearchWisataRating(w http.ResponseWriter, r *http.Request) {
	// Ambil query parameter 'search' dari URL
	searchQuery := r.URL.Query().Get("search")
	if searchQuery == "" {
		// Jika query search kosong, kirimkan error
		http.Error(w, "Search query is required", http.StatusBadRequest)
		return
	}

	// Inisialisasi slice untuk menampung hasil pencarian
	var wisatas []models.Wisata

	// Query pencarian berdasarkan 'nama_wisata', mencari yang mengandung kata searchQuery
	// Menggunakan ILIKE untuk pencarian case-insensitive di PostgreSQL
	if result := config.DB.Where("nama_wisata ILIKE ?", "%"+searchQuery+"%").Find(&wisatas); result.Error != nil {
		// Jika ada error dalam query
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Menampung hasil pencarian dengan rating rata-rata dan jumlah pemberi rating
	var wisataWithRating []map[string]interface{}
	for _, w := range wisatas {
		// Mengambil data rating untuk wisata
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

		// Menambahkan data wisata, rating rata-rata, dan jumlah pemberi rating
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

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wisataWithRating)
}

// Fungsi untuk mendapatkan detail wisata berdasarkan id, dengan rata-rata rating, jumlah pemberi rating dan ulasan
func GetWisataDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data wisata berdasarkan id
	var wisata models.Wisata
	if err := config.DB.Where("id_wisata = ?", id).First(&wisata).Error; err != nil {
		http.Error(w, "Data wisata not found", http.StatusNotFound)
		return
	}

	// Kirimkan detail data wisata terlebih dahulu
	response := map[string]interface{}{
		"wisata": wisata,
	}

	// Ambil data rating untuk wisata
	var rating []models.RatingWisata
	if err := config.DB.Where("id_wisata = ?", wisata.IDWisata).Find(&rating).Error; err != nil {
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
	config.DB.Model(&models.RatingWisata{}).Where("id_wisata = ?", wisata.IDWisata).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response["rating_average"] = averageRating
	response["rating_counts"] = ratingCounts
	response["reviewers_count"] = uniqueRaters
	response["ratings"] = rating

	// Menambahkan data wisata lainnya
	var wisataLain []models.Wisata
	if err := config.DB.Where("id_wisata != ?", id).Find(&wisataLain).Error; err != nil {
		http.Error(w, "Data wisata lain not found", http.StatusInternalServerError)
		return
	}

	var wisataLainWithDetails []map[string]interface{}
	for _, w := range wisataLain {
		var wisataData = map[string]interface{}{
			"id_wisata":             w.IDWisata,
			"nama_wisata":           w.NamaWisata,
			"lokasi_wisata":         w.LokasiWisata,
			"koordinat_wisata":      w.KoordinatWisata,
			"kontak_wisata":         w.KontakWisata,
			"foto_wisata":           w.FotoWisata,
			"deskripsi":             w.Deskripsi,
			"fasilitas":             w.Fasilitas,
			"rating":                averageRating, // Set default average rating (for demo)
			"jumlah_pemberi_rating": uniqueRaters,  // Set default reviewer count
		}
		wisataLainWithDetails = append(wisataLainWithDetails, wisataData)
	}

	// Menambahkan data wisata lainnya ke dalam response
	response["wisata_lainnnya"] = wisataLainWithDetails

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Fungsi untuk mendapatkan detail rating wisata berdasarkan id, dengan rata-rata rating, jumlah pemberi rating, dan ulasan
func GetWisataRatingDetailByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Ambil data wisata berdasarkan id
	var wisata models.Wisata
	if err := config.DB.Where("id_wisata = ?", id).First(&wisata).Error; err != nil {
		http.Error(w, "Data wisata not found", http.StatusNotFound)
		return
	}

	// Kirimkan detail data wisata terlebih dahulu
	response := map[string]interface{}{
		"wisata": wisata,
	}

	// Ambil data rating untuk wisata
	var rating []models.RatingWisata
	if err := config.DB.Where("id_wisata = ?", wisata.IDWisata).Find(&rating).Error; err != nil {
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
	config.DB.Model(&models.RatingWisata{}).Where("id_wisata = ?", wisata.IDWisata).Distinct("id_akun").Count(&uniqueRaters)

	// Menambahkan data rating rata-rata, jumlah pemberi rating dan rating counts
	response["rating_average"] = averageRating
	response["rating_counts"] = ratingCounts
	response["reviewers_count"] = uniqueRaters
	response["ratings"] = rating

	// Kirimkan response ke client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
