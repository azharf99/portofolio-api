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
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if industry != "" {
		query = query.Where("industry = ?", industry)
	}
	if pType != "" {
		query = query.Where("type = ?", pType)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Find(&portfolios).Error
	return portfolios, total, err
}

func (r *portfolioRepository) Store(portfolio *domain.Portfolio) error {
	return r.db.Create(portfolio).Error
}

func (r *portfolioRepository) Update(id uint, portfolio *domain.Portfolio) error {
	// Updates() akan memperbarui field yang tidak kosong saja
	return r.db.Model(&domain.Portfolio{}).Where("id = ?", id).Updates(portfolio).Error
}

func (r *portfolioRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Portfolio{}, id).Error
}
