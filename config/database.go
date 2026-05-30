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
	err = db.AutoMigrate(&domain.Portfolio{}, &domain.User{}, &domain.PortfolioImage{}, &domain.Service{}, &domain.Transaction{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	seedAdmin(db)
	seedServices(db)

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

func seedServices(db *gorm.DB) {
	var count int64
	db.Model(&domain.Service{}).Count(&count)

	if count == 0 {
		services := []domain.Service{
			{
				Title:         "Landing Page Development",
				Description:   "Create a premium, high-converting responsive landing page tailored for your business needs.",
				OriginalPrice: 500000,
				PromoPrice:    300000,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20Landing%20Page%20Development%20service",
				IsActive:      true,
			},
			{
				Title:         "Point of Sales (POS) System",
				Description:   "A complete retail/restaurant POS system with inventory tracking, sales reports, and receipt printing.",
				OriginalPrice: 4500000,
				PromoPrice:    3000000,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20POS%20System%20service",
				IsActive:      true,
			},
			{
				Title:         "Enterprise Resource Planning (ERP) Suite",
				Description:   "Full-scale ERP system including HR, finance, supply chain, and customer relationship management.",
				OriginalPrice: 8000000,
				PromoPrice:    5000000,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20ERP%20Suite%20service",
				IsActive:      true,
			},
		}
		for _, s := range services {
			db.Create(&s)
		}
		log.Println("Initial services data seeded successfully!")
	}
}

