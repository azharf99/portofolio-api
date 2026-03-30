package mocks

import (
	"github.com/azharf99/portofolio-api/domain"
	"github.com/stretchr/testify/mock"
)

type UserUsecaseMock struct {
	mock.Mock
}

func (m *UserUsecaseMock) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *UserUsecaseMock) Update(id uint, user *domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *UserUsecaseMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
