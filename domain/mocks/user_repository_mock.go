package mocks

import (
	"github.com/azharf99/portofolio-api/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetByUsername(username string) (domain.User, error) {
	args := m.Called(username)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByID(id uint) (domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(domain.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(id uint, user *domain.User) error {
	args := m.Called(id, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
