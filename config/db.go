package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB adalah koneksi global ke database PostgreSQL
var DB *gorm.DB

// ConnectDatabase menghubungkan aplikasi Go dengan database PostgreSQL di Google Cloud
func ConnectDatabase() {
	// Ganti dengan data kredensial Google Cloud SQL Anda
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

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
