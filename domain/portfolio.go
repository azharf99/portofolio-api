package domain

import "time"

type Portfolio struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Role        string    `json:"role"`
	Type        string    `json:"type"`
	Industry    string    `json:"industry"`
	TechStack   string    `json:"tech_stack"`
	ProjectLink string    `json:"project_link"`
	ImageURL    string    `json:"image_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsPublished bool      `json:"is_published" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PortfolioUsecase interface {
	Fetch(page, limit int, search, industry, pType string) ([]Portfolio, int64, error)
	Store(portfolio *Portfolio) error
	Update(id uint, portfolio *Portfolio) error // BARU
	Delete(id uint) error                       // BARU
}

type PortfolioRepository interface {
	Fetch(limit, offset int, search, industry, pType string) ([]Portfolio, int64, error)
	Store(portfolio *Portfolio) error
	Update(id uint, portfolio *Portfolio) error // BARU
	Delete(id uint) error                       // BARU
}
