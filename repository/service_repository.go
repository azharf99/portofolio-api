package repository

import (
	"github.com/azharf99/portofolio-api/domain"
	"gorm.io/gorm"
)

type serviceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) domain.ServiceRepository {
	return &serviceRepository{db}
}

func (r *serviceRepository) Fetch(limit, offset int, activeOnly bool) ([]domain.Service, int64, error) {
	var services []domain.Service
	var total int64

	query := r.db.Model(&domain.Service{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	query.Count(&total)
	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&services).Error
	return services, total, err
}

func (r *serviceRepository) GetByID(id uint) (*domain.Service, error) {
	var service domain.Service
	err := r.db.First(&service, id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) Store(service *domain.Service) error {
	return r.db.Create(service).Error
}

func (r *serviceRepository) Update(id uint, service *domain.Service) error {
	var existing domain.Service
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}

	service.ID = id
	return r.db.Save(service).Error
}

func (r *serviceRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Service{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
