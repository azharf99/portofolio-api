package usecase

import "github.com/azharf99/portofolio-api/domain"

type portfolioUsecase struct {
	repo domain.PortfolioRepository
}

func NewPortfolioUsecase(repo domain.PortfolioRepository) domain.PortfolioUsecase {
	return &portfolioUsecase{repo}
}

func (u *portfolioUsecase) Fetch(page, limit int, search, industry, pType string, onlyPublished bool) ([]domain.Portfolio, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	return u.repo.Fetch(limit, offset, search, industry, pType, onlyPublished)
}

func (u *portfolioUsecase) Store(portfolio *domain.Portfolio) error {
	// Di sini kamu bisa menambahkan validasi, misalnya mengecek URL valid atau tidak
	return u.repo.Store(portfolio)
}

func (u *portfolioUsecase) Update(id uint, portfolio *domain.Portfolio) error {
	return u.repo.Update(id, portfolio)
}

func (u *portfolioUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}
