package routes

import (
	"log"
	"os"

	"github.com/azharf99/portofolio-api/delivery/http"
	"github.com/azharf99/portofolio-api/middleware"
	"github.com/azharf99/portofolio-api/repository"
	"github.com/azharf99/portofolio-api/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes mengatur semua inisialisasi layer dan rute API
func SetupRoutes(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	// Buat folder upload jika belum ada
	_ = os.MkdirAll("uploads/portfolios", os.ModePerm)
	_ = os.MkdirAll("uploads/services", os.ModePerm)

	// 1. Setup Repository
	portfolioRepo := repository.NewPortfolioRepository(db)
	userRepo := repository.NewUserRepository(db)
	serviceRepo := repository.NewServiceRepository(db)

	// 2. Setup Usecase
	portfolioUsecase := usecase.NewPortfolioUsecase(portfolioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo)

	// 3. Jalankan Cleanup Gambar yang hilang di disk (Opsional tapi berguna jika kontainer restart tanpa volume)
	if err := portfolioUsecase.CleanupOrphanedImages(); err != nil {
		log.Printf("Warning: Gagal melakukan cleanup gambar: %v\n", err)
	}

	// 3. Setup Handler (Controller)
	portfolioHandler := http.NewPortfolioHandlerInstance(portfolioUsecase)
	userHandler := http.NewUserHandlerInstance(userUsecase)
	serviceHandler := http.NewServiceHandlerInstance(serviceUsecase)

	// 4. Daftarkan Rute
	api := r.Group("/api")

	// === PUBLIC ROUTES ===
	r.Static("/uploads", "./uploads")
	api.POST("/login", userHandler.Login)
	api.GET("/portfolios", portfolioHandler.Fetch)
	api.GET("/services", serviceHandler.Fetch)

	// === PRIVATE ROUTES (Butuh Login) ===
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret))

	// CRUD Portofolio Private
	admin.GET("/portfolios", portfolioHandler.AdminFetch)
	admin.POST("/portfolios", portfolioHandler.Store)
	admin.PUT("/portfolios/:id", portfolioHandler.Update)
	admin.DELETE("/portfolios/:id", portfolioHandler.Delete)

	// CRUD Service Private
	admin.GET("/services", serviceHandler.AdminFetch)
	admin.POST("/services", serviceHandler.Store)
	admin.PUT("/services/:id", serviceHandler.Update)
	admin.DELETE("/services/:id", serviceHandler.Delete)

	// Manajemen User Private
	admin.PUT("/users/:id", userHandler.Update)
	admin.DELETE("/users/:id", userHandler.Delete)
}

