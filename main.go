package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/azharf99/portofolio-api/config"
	"github.com/azharf99/portofolio-api/middleware"
	"github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/azharf99/portofolio-api/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup Database and Env
	db := config.SetupDatabase()
	i18n.Init() // Initialize i18n

	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan. Membaca konfigurasi dari Docker Environment.")
	}

	// Ambil JWT Secret dari env, fallback jika kosong
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET wajib diisi untuk menjalankan aplikasi")
	}

	// 2. Setup Framework Gin
	r := gin.Default()

	// KEAMANAN: Gin secara default mempercayai header X-Forwarded-For/X-Real-IP dari
	// SIAPA SAJA jika SetTrustedProxies tidak diset. Ini membuat c.ClientIP() (dipakai
	// oleh RateLimiter) mudah dipalsukan (spoofed) untuk melewati rate limiting.
	// Set TRUSTED_PROXIES=ip1,ip2 di env jika API ini berada di belakang reverse proxy
	// (nginx/Traefik/Cloudflare). Default: jangan percaya proxy manapun, pakai IP koneksi asli.
	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		if err := r.SetTrustedProxies(nil); err != nil {
			log.Fatal("Gagal mengatur trusted proxies:", err)
		}
	} else {
		if err := r.SetTrustedProxies(strings.Split(trustedProxies, ",")); err != nil {
			log.Fatal("Gagal mengatur trusted proxies:", err)
		}
	}

	// 3. Pasang Middleware Keamanan Global & CORS
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.SetupCORS())
	r.Use(middleware.RateLimiter())
	r.Use(middleware.I18nMiddleware())

	// 4. Setup Routing (Memanggil dari package routes)
	routes.SetupRoutes(r, db, jwtSecret)

	// 5. Jalankan Aplikasi
	fmt.Println("🚀 Portofolio API berjalan di port 8080")
	r.Run(":8080")
}
