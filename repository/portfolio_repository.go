package repository

import (
	"log"
	"os"
	"strings"

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
	err := query.Preload("Images").Offset(offset).Limit(limit).Order("end_date DESC").Find(&portfolios).Error
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

func (r *portfolioRepository) CleanupOrphanedImages() error {
	// 1. Cek Main Image di tabel Portfolios
	var portfolios []domain.Portfolio
	if err := r.db.Find(&portfolios).Error; err != nil {
		return err
	}

	for _, p := range portfolios {
		if p.ImageURL != "" {
			// Hilangkan leading slash jika ada (misal /uploads/... -> uploads/...)
			filePath := strings.TrimPrefix(p.ImageURL, "/")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				log.Printf("Cleaning up missing main image: %s for Portfolio ID: %d\n", p.ImageURL, p.ID)
				r.db.Model(&p).Update("image_url", "")
			}
		}
	}

	// 2. Cek Gallery Images di tabel PortfolioImages
	var galleryImages []domain.PortfolioImage
	if err := r.db.Find(&galleryImages).Error; err != nil {
		return err
	}

	for _, img := range galleryImages {
		if img.ImageURL != "" {
			filePath := strings.TrimPrefix(img.ImageURL, "/")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				log.Printf("Deleting missing gallery image record: %s (ID: %d)\n", img.ImageURL, img.ID)
				r.db.Delete(&img)
			}
		}
	}

	return nil
}
