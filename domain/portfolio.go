package domain

import "time"

type Portfolio struct {
	ID          uint             `json:"id" gorm:"primaryKey" form:"id"`
	Title       string           `json:"title" form:"title"`
	Description string           `json:"description" form:"description"`
	Role        string           `json:"role" form:"role"`
	Type        string           `json:"type" form:"type"`
	Industry    string           `json:"industry" form:"industry"`
	TechStack   string           `json:"tech_stack" form:"tech_stack"`
	ProjectLink string           `json:"project_link" form:"project_link"`
	ImageURL    string           `json:"image_url" form:"image_url"` // Main thumbnail
	Images      []PortfolioImage `json:"images" gorm:"foreignKey:PortfolioID;constraint:OnDelete:CASCADE"`
	StartDate   time.Time        `json:"start_date" form:"start_date" time_format:"2006-01-02"`
	EndDate     time.Time        `json:"end_date" form:"end_date" time_format:"2006-01-02"`
	IsPublished bool             `json:"is_published" gorm:"default:true" form:"is_published"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type PortfolioImage struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PortfolioID uint      `json:"portfolio_id"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PortfolioUsecase interface {
	Fetch(page, limit int, search, industry, pType string, onlyPublished bool) ([]Portfolio, int64, error)
	Store(portfolio *Portfolio) error
	Update(id uint, portfolio *Portfolio) error // BARU
	Delete(id uint) error                       // BARU
}

type PortfolioRepository interface {
	Fetch(limit, offset int, search, industry, pType string, onlyPublished bool) ([]Portfolio, int64, error)
	Store(portfolio *Portfolio) error
	Update(id uint, portfolio *Portfolio) error // BARU
	Delete(id uint) error                       // BARU
}
