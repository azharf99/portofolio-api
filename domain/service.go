package domain

import "time"

type Service struct {
	ID            uint      `json:"id" gorm:"primaryKey" form:"id"`
	Title         string    `json:"title" form:"title"`
	Description   string    `json:"description" form:"description"`
	OriginalPrice int64     `json:"original_price" form:"original_price"`
	PromoPrice    int64     `json:"promo_price" form:"promo_price"`
	ImageURL      string    `json:"image_url" form:"image_url"`
	RedirectURL   string    `json:"redirect_url" form:"redirect_url"`
	IsActive      bool      `json:"is_active" gorm:"default:true" form:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ServiceUsecase interface {
	Fetch(page, limit int, activeOnly bool) ([]Service, int64, error)
	GetByID(id uint) (*Service, error)
	Store(service *Service) error
	Update(id uint, service *Service) error
	Delete(id uint) error
}

type ServiceRepository interface {
	Fetch(limit, offset int, activeOnly bool) ([]Service, int64, error)
	GetByID(id uint) (*Service, error)
	Store(service *Service) error
	Update(id uint, service *Service) error
	Delete(id uint) error
}
