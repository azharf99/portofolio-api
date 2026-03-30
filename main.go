package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/azharf99/portofolio-api/config"
	"github.com/azharf99/portofolio-api/middleware"
	"github.com/azharf99/portofolio-api/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Setup Database and Env
	db := config.SetupDatabase()

	err := godotenv.Load()
	if err != nil {
		log.Println("Info: File .env tidak ditemukan. Membaca konfigurasi dari Docker Environment.")
	}

	// Ambil JWT Secret dari env, fallback jika kosong
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "RAHASIA_NEGARA_SANGAT_AMAN_SEKALI_99"
	}

	// 2. Setup Framework Gin
	r := gin.Default()

	// 3. Pasang Middleware Keamanan Global & CORS
	r.Use(middleware.SecurityHeaders())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Ganti dengan domainmu nanti saat deploy
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 4. Setup Routing (Memanggil dari package routes)
	routes.SetupRoutes(r, db, jwtSecret)

	// 5. Jalankan Aplikasi
	fmt.Println("🚀 Portofolio API berjalan di port 8080")
	r.Run(":8080")
}
