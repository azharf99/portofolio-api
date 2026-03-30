package config

import (
	"fmt"
	"log"
	"os"

	"github.com/azharf99/portofolio-api/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupDatabase bertugas membuka koneksi ke PostgreSQL
func SetupDatabase() *gorm.DB {
	// Membaca konfigurasi dari Environment Variables (diset via Docker nanti)

	// ==========================================
	// 2. KONFIGURASI DATABASE POSTGRESQL
	// ==========================================
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database PostgreSQL:", err)
	}
	fmt.Println("Berhasil terhubung ke database PostgreSQL!")

	// Migrasi tabel otomatis
	err = db.AutoMigrate(&domain.Portfolio{}, &domain.User{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	seedAdmin(db)

	return db
}

func seedAdmin(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Count(&count)

	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := domain.User{
			Username: "azharfa",
			Password: string(hashedPassword),
		}
		db.Create(&admin)
		log.Println("User admin berhasil dibuat! Username: azharfa | Password: admin123")
	}
}
