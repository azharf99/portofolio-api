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

func (r *portfolioRepository) Fetch(limit, offset int, search, industry, pType string) ([]domain.Portfolio, int64, error) {
	var portfolios []domain.Portfolio
	var total int64

	query := r.db.Model(&domain.Portfolio{}).Where("is_published = ?", true)

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
	err := query.Offset(offset).Limit(limit).Find(&portfolios).Error
	return portfolios, total, err
}

func (r *portfolioRepository) Store(portfolio *domain.Portfolio) error {
	return r.db.Create(portfolio).Error
}

func (r *portfolioRepository) Update(id uint, portfolio *domain.Portfolio) error {
	result := r.db.Model(&domain.Portfolio{}).Where("id = ?", id).Updates(portfolio)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
