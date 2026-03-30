package mocks

import (
	"github.com/azharf99/portofolio-api/domain"
	"github.com/stretchr/testify/mock"
)

type PortfolioRepositoryMock struct {
	mock.Mock
}

func (m *PortfolioRepositoryMock) Fetch(limit, offset int, search, industry, pType string) ([]domain.Portfolio, int64, error) {
	args := m.Called(limit, offset, search, industry, pType)
	return args.Get(0).([]domain.Portfolio), args.Get(1).(int64), args.Error(2)
}

func (m *PortfolioRepositoryMock) Store(portfolio *domain.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *PortfolioRepositoryMock) Update(id uint, portfolio *domain.Portfolio) error {
	args := m.Called(id, portfolio)
	return args.Error(0)
}

func (m *PortfolioRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
