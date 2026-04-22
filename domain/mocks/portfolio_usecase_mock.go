package mocks

import (
	"github.com/azharf99/portofolio-api/domain"
	"github.com/stretchr/testify/mock"
)

type PortfolioUsecaseMock struct {
	mock.Mock
}

func (m *PortfolioUsecaseMock) Fetch(page, limit int, search, industry, pType string, onlyPublished bool) ([]domain.Portfolio, int64, error) {
	args := m.Called(page, limit, search, industry, pType, onlyPublished)
	return args.Get(0).([]domain.Portfolio), args.Get(1).(int64), args.Error(2)
}

func (m *PortfolioUsecaseMock) Store(portfolio *domain.Portfolio) error {
	args := m.Called(portfolio)
	return args.Error(0)
}

func (m *PortfolioUsecaseMock) Update(id uint, portfolio *domain.Portfolio) error {
	args := m.Called(id, portfolio)
	return args.Error(0)
}

func (m *PortfolioUsecaseMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}