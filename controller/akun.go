package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"
	"xplore/config"
	"xplore/models"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	var akun models.Akun
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&akun); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mencari akun admin berdasarkan username, password, dan user_role = 'admin'
	if err := config.DB.Where("username = ? AND password = ? AND user_role = ?", akun.Username, akun.Password, "admin").First(&akun).Error; err != nil {
		http.Error(w, "User not found or incorrect credentials", http.StatusUnauthorized)
		return
	}

	// Mengembalikan seluruh data akun yang ditemukan
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun) // Mengirimkan seluruh objek akun
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var akun models.Akun
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&akun); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Mencari akun user berdasarkan email, password, dan user_role = 'user'
	if err := config.DB.Where("email = ? AND password = ? AND user_role = ?", akun.Email, akun.Password, "user").First(&akun).Error; err != nil {
		http.Error(w, "User not found or incorrect credentials", http.StatusUnauthorized)
		return
	}

	// Mengembalikan seluruh data akun yang ditemukan
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun) // Mengirimkan seluruh objek akun
}

// Fungsi untuk mengupload foto ke Google Cloud Storage
func uploadFotoToGCS(file multipart.File, filename string) (string, error) {
	by, err := json.Marshal(config.Config)
	if err != nil {
		return "", fmt.Errorf("error marshall : %v", err.Error())
	}

	// Membuat klien Google Cloud Storage
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(by))
	if err != nil {
		return "", fmt.Errorf("failed to create GCS client: %v", err)
	}
	defer client.Close()

	// Tentukan bucket dan objek
	bucket := client.Bucket("xplores") // Ganti dengan nama bucket Anda
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

// Fungsi untuk memperbarui profil akun admin
func UpdateAdminAkun(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil ID Akun dari URL params
	vars := mux.Vars(r)
	idAkun := vars["id"]

	// Ambil data dari form
	username := r.FormValue("username")
	nohp := r.FormValue("nohp")
	email := r.FormValue("email")
	password := r.FormValue("password")
	userRole := "admin" // Set userRole secara otomatis menjadi "admin"

	// Ambil file foto profil baru
	files := r.MultipartForm.File["foto_profil"]
	var fotoProfilLinks []string

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
		link, err := uploadFotoToGCS(file, fileHeader.Filename)
		if err != nil {
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoProfilLinks = append(fotoProfilLinks, link)
	}

	// Ambil data akun yang ada di database
	var akun models.Akun
	if err := config.DB.First(&akun, idAkun).Error; err != nil {
		http.Error(w, "Akun not found", http.StatusNotFound)
		return
	}

	// Update data akun
	akun.Username = username
	akun.Nohp = nohp
	akun.Email = email
	akun.Password = password
	akun.UserRole = userRole

	// Tambahkan foto profil baru jika ada, namun foto lama tetap ada
	if len(fotoProfilLinks) > 0 {
		akun.Fotoprofil = fotoProfilLinks[0] // Gantikan foto lama dengan foto baru
	}

	// Update timestamp
	akun.UpdatedAt = time.Now()

	// Simpan update ke database
	if err := config.DB.Save(&akun).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun)
}

// GetAllAkun untuk mendapatkan semua akun
func GetAllAkun(w http.ResponseWriter, r *http.Request) {
	var akuns []models.Akun
	if result := config.DB.Find(&akuns); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akuns)
}

// GetAkunByID untuk mendapatkan akun berdasarkan ID
func GetAkunByID(w http.ResponseWriter, r *http.Request) {
	var akun models.Akun
	id := mux.Vars(r)["id"]

	// Menyaring akun berdasarkan ID
	if result := config.DB.First(&akun, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun)
}

// Fungsi untuk memperbarui profil akun (user) dan mengupload foto profil baru
func UpdateAkun(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB limit
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Ambil ID Akun dari URL params
	vars := mux.Vars(r)
	idAkun := vars["id"]

	// Ambil data lainnya dari form
	username := r.FormValue("username")
	nohp := r.FormValue("nohp")
	email := r.FormValue("email")
	password := r.FormValue("password")
	userRole := "user" // Set userRole secara otomatis menjadi "user"

	// Ambil file foto profil baru
	files := r.MultipartForm.File["foto_profil"]
	var fotoProfilLinks []string

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
		link, err := uploadFotoToGCS(file, fileHeader.Filename)
		if err != nil {
			log.Printf("Error when upload photo: %v", err.Error())
			http.Error(w, "Unable to upload to GCS", http.StatusInternalServerError)
			return
		}

		// Simpan URL file di array
		fotoProfilLinks = append(fotoProfilLinks, link)
	}

	// Ambil data akun yang ada di database
	var akun models.Akun
	if err := config.DB.First(&akun, idAkun).Error; err != nil {
		http.Error(w, "Akun not found", http.StatusNotFound)
		return
	}

	// Update data akun
	akun.Username = username
	akun.Nohp = nohp
	akun.Email = email
	akun.Password = password
	akun.UserRole = userRole

	// Jika ada foto baru, perbarui foto profilnya dengan array foto baru
	if len(fotoProfilLinks) > 0 {
		akun.Fotoprofil = fotoProfilLinks[0] // Gantikan foto lama dengan foto baru
	}

	// Update timestamp
	akun.UpdatedAt = time.Now()

	// Simpan update ke database
	if err := config.DB.Save(&akun).Error; err != nil {
		http.Error(w, "Unable to save to database", http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun)
}

// DeleteAkun untuk menghapus akun berdasarkan ID
func DeleteAkun(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if result := config.DB.Delete(&models.Akun{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Fungsi untuk membuat akun baru
func CreateAkun(w http.ResponseWriter, r *http.Request) {
	var akun models.Akun

	// Decode JSON body request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&akun); err != nil {
		http.Error(w, "Unable to parse request", http.StatusBadRequest)
		return
	}

	// Validasi email dan password
	if akun.Email == "" || akun.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	// Mengisi field lain dengan template yang sudah ditentukan
	akun.Username = "user123"
	akun.Nohp = "081234567890"
	akun.Fotoprofil = "profile_pic.jpg"
	akun.UserRole = "user"

	// Set CreatedAt dan UpdatedAt
	akun.CreatedAt = time.Now()
	akun.UpdatedAt = time.Now()

	// Simpan akun ke database
	if err := config.DB.Create(&akun).Error; err != nil {
		http.Error(w, fmt.Sprintf("Failed to create account: %v", err), http.StatusInternalServerError)
		return
	}

	// Kirim respons sukses
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun)
}

// Fungsi untuk mencari akun berdasarkan email
func GetAkunByEmail(w http.ResponseWriter, r *http.Request) {
	// Ambil email dari URL params
	vars := mux.Vars(r)
	email := vars["email"]

	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Mencari akun berdasarkan email
	var akun models.Akun
	if err := config.DB.Where("email = ?", email).First(&akun).Error; err != nil {
		http.Error(w, fmt.Sprintf("Akun dengan email %s tidak ditemukan", email), http.StatusNotFound)
		return
	}

	// Kirimkan data akun yang ditemukan
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(akun)
}
