package usecase

import "github.com/azharf99/portofolio-api/domain"

type serviceUsecase struct {
	repo domain.ServiceRepository
}

func NewServiceUsecase(repo domain.ServiceRepository) domain.ServiceUsecase {
	return &serviceUsecase{repo}
}

func (u *serviceUsecase) Fetch(page, limit int, activeOnly bool) ([]domain.Service, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return u.repo.Fetch(limit, offset, activeOnly)
}

func (u *serviceUsecase) Store(service *domain.Service) error {
	return u.repo.Store(service)
}

func (u *serviceUsecase) Update(id uint, service *domain.Service) error {
	return u.repo.Update(id, service)
}

func (u *serviceUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
