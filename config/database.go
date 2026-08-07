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
		// KEAMANAN: Jangan pernah hardcode username/password default (mis. admin123).
		// Wajib diisi eksplisit lewat env, sama seperti JWT_SECRET, agar operator
		// sadar membuat kredensial sendiri sebelum API ini live.
		username := os.Getenv("ADMIN_USERNAME")
		password := os.Getenv("ADMIN_PASSWORD")
		if username == "" || password == "" {
			log.Fatal("ADMIN_USERNAME dan ADMIN_PASSWORD wajib diisi di env untuk membuat akun admin pertama kali")
		}
		if len(password) < 12 {
			log.Fatal("ADMIN_PASSWORD wajib minimal 12 karakter")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Gagal melakukan hashing password admin:", err)
		}
		admin := domain.User{
			Username: username,
			Password: string(hashedPassword),
		}
		db.Create(&admin)
		// KEAMANAN: Jangan log password, bahkan yang sudah di-hash, ke stdout/aggregator log.
		log.Printf("User admin '%s' berhasil dibuat.\n", username)
	}
}

// seedServices membuat 3 paket bundle (bukan lagi item satuan) — dikelompokkan
// per jenis bisnis sesuai arah baru: UMKM (kafe/resto), Sekolah, dan Enterprise.
// Setiap bundle menyertakan WhatsApp Gateway sebagai fitur default.
// Paket Sekolah sengaja diberi harga 0 (custom quote) — frontend menampilkan
// "Custom quote" dan mengarahkan CTA ke RedirectURL (WhatsApp) alih-alih ke
// alur checkout/pembayaran ketika kedua harga bernilai 0.
func seedServices(db *gorm.DB) {
	var count int64
	db.Model(&domain.Service{}).Count(&count)

	if count == 0 {
		services := []domain.Service{
			{
				Title:         "Paket UMKM — Kafe & Restoran",
				Description:   "Sistem POS lengkap dengan manajemen stok dan laporan penjualan, landing page 1 halaman untuk kehadiran online, dan integrasi WhatsApp Gateway untuk notifikasi pesanan otomatis — siap dipakai sejak hari pertama.",
				Features:      "Sistem POS dengan manajemen stok,Cetak struk otomatis,Landing page bisnis 1 halaman,Integrasi WhatsApp Gateway,Laporan penjualan harian",
				OriginalPrice: 4500000,
				PromoPrice:    3000000,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20UMKM%20Starter%20package",
				IsActive:      true,
			},
			{
				Title:         "Paket Sekolah & Lembaga Kursus",
				Description:   "Learning Management System untuk pengelolaan kelas, tugas, dan nilai, dilengkapi portal siswa/orang tua serta integrasi WhatsApp Gateway untuk notifikasi otomatis. Harga menyesuaikan jumlah siswa dan fitur yang dibutuhkan.",
				Features:      "Manajemen kelas & kohort,Tugas & penilaian online,Portal siswa/orang tua,Integrasi WhatsApp Gateway,Skalabel sesuai jumlah siswa",
				OriginalPrice: 0,
				PromoPrice:    0,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20Learning%20Suite%20package%20for%20my%20school",
				IsActive:      true,
			},
			{
				Title:         "Paket Enterprise — ERP Suite",
				Description:   "ERP terpadu mencakup HR & payroll, keuangan, supply chain, dan CRM dalam satu sistem, dengan integrasi WhatsApp Gateway untuk notifikasi lintas tim — menggantikan kombinasi banyak spreadsheet dengan satu sumber data yang konsisten.",
				Features:      "Modul HR & payroll,Keuangan & supply chain,CRM & dashboard laporan,Integrasi WhatsApp Gateway,Multi-user & manajemen akses",
				OriginalPrice: 15000000,
				PromoPrice:    10000000,
				ImageURL:      "",
				RedirectURL:   "https://wa.me/6285702570200?text=Hello%2C%20I%20am%20interested%20in%20your%20Enterprise%20Suite%20package",
				IsActive:      true,
			},
		}
		for _, s := range services {
			db.Create(&s)
		}
		log.Println("Initial bundle services seeded successfully!")
	}
}
