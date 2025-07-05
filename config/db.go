package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah koneksi global ke database PostgreSQL
var DB *gorm.DB

// ConnectDatabase menghubungkan aplikasi Go dengan database PostgreSQL di Google Cloud
func ConnectDatabase() {
	// Ganti dengan data kredensial Google Cloud SQL Anda
	host := "xplore-48:asia-southeast2:xplore-db"
	user := "xplore-db"       // Username database
	password := "damar150705" // Password database
	dbname := "xplore_db"     // Nama database
	port := "5432"            // Port PostgreSQL (default)

	// Membuat string koneksi
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	// Menghubungkan ke database menggunakan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	// Menyimpan koneksi global untuk digunakan di seluruh aplikasi
	DB = database
	fmt.Println("Koneksi ke database berhasil!")
}
