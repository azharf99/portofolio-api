package repository

import (
	"github.com/azharf99/portofolio-api/domain"
	"gorm.io/gorm"
)

type portfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository membuat instance baru dari repository portofolio
func NewPortfolioRepository(db *gorm.DB) domain.PortfolioRepository {
	return &portfolioRepository{db}
}

func (r *portfolioRepository) Fetch(limit, offset int, search, industry, pType string, onlyPublished bool) ([]domain.Portfolio, int64, error) {
	var portfolios []domain.Portfolio
	var total int64

	query := r.db.Model(&domain.Portfolio{})

	if onlyPublished {
		query = query.Where("is_published = ?", true)
	}

	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if industry != "" {
		query = query.Where("industry ILIKE ?", "%"+industry+"%")
	}
	if pType != "" {
		query = query.Where("type ILIKE ?", "%"+pType+"%")
	}

	query.Count(&total)
	err := query.Preload("Images").Offset(offset).Limit(limit).Order("created_at DESC").Find(&portfolios).Error
	return portfolios, total, err
}

func (r *portfolioRepository) Store(portfolio *domain.Portfolio) error {
	return r.db.Create(portfolio).Error
}

func (r *portfolioRepository) Update(id uint, portfolio *domain.Portfolio) error {
	// 1. Cek apakah record ada
	var existing domain.Portfolio
	if err := r.db.First(&existing, id).Error; err != nil {
		return err // gorm.ErrRecordNotFound if not found
	}

	// 2. Lakukan Update
	// Kita gunakan Save untuk menyinkronkan seluruh field termasuk asosiasi Images
	portfolio.ID = id
	return r.db.Save(portfolio).Error
}

func (r *portfolioRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Portfolio{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
