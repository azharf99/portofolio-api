package routes

import (
	"github.com/azharf99/portofolio-api/delivery/http"
	"github.com/azharf99/portofolio-api/middleware"
	"github.com/azharf99/portofolio-api/repository"
	"github.com/azharf99/portofolio-api/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes mengatur semua inisialisasi layer dan rute API
func SetupRoutes(r *gin.Engine, db *gorm.DB, jwtSecret string) {
	// 1. Setup Repository
	portfolioRepo := repository.NewPortfolioRepository(db)
	userRepo := repository.NewUserRepository(db)

	// 2. Setup Usecase
	portfolioUsecase := usecase.NewPortfolioUsecase(portfolioRepo)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)

	// 3. Setup Handler (Controller)
	portfolioHandler := http.NewPortfolioHandlerInstance(portfolioUsecase)
	userHandler := http.NewUserHandlerInstance(userUsecase)

	// 4. Daftarkan Rute
	api := r.Group("/api")

	// === PUBLIC ROUTES ===
	r.Static("/uploads", "./uploads")
	api.POST("/login", userHandler.Login)
	api.GET("/portfolios", portfolioHandler.Fetch)

	// === PRIVATE ROUTES (Butuh Login) ===
	admin := api.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret))

	// CRUD Portofolio Private
	admin.GET("/portfolios", portfolioHandler.AdminFetch)
	admin.POST("/portfolios", portfolioHandler.Store)
	admin.PUT("/portfolios/:id", portfolioHandler.Update)
	admin.DELETE("/portfolios/:id", portfolioHandler.Delete)

	// Manajemen User Private
	admin.PUT("/users/:id", userHandler.Update)
	admin.DELETE("/users/:id", userHandler.Delete)
}
