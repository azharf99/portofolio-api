package usecase_test

import (
	"testing"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/azharf99/portofolio-api/domain/mocks"
	"github.com/azharf99/portofolio-api/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserUsecase_Login(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	jwtSecret := "secret_test"
	mockUsecase := usecase.NewUserUsecase(mockRepo, jwtSecret)

	t.Run("Success Login", func(t *testing.T) {
		// Kita harus membuat hash password asli agar lolos pengecekan bcrypt di usecase
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		mockUser := domain.User{
			ID:       1,
			Username: "azharfa",
			Password: string(hashedPassword),
		}

		// Ekspektasi: Repository mengembalikan data user beserta password yang di-hash
		mockRepo.On("GetByUsername", "azharfa").Return(mockUser, nil).Once()

		token, err := mockUsecase.Login("azharfa", "admin123")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Failed Login - Wrong Password", func(t *testing.T) {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		mockUser := domain.User{
			ID:       1,
			Username: "azharfa",
			Password: string(hashedPassword),
		}

		mockRepo.On("GetByUsername", "azharfa").Return(mockUser, nil).Once()

		// Kita coba login dengan password yang salah
		token, err := mockUsecase.Login("azharfa", "salahpassword")

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Equal(t, "username atau password salah", err.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Update(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockUsecase := usecase.NewUserUsecase(mockRepo, "secret")

	t.Run("Success Update", func(t *testing.T) {
		mockUser := &domain.User{Username: "azharfa_baru"}

		mockRepo.On("Update", uint(1), mock.AnythingOfType("*domain.User")).Return(nil).Once()

		err := mockUsecase.Update(1, mockUser)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserUsecase_Delete(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	mockUsecase := usecase.NewUserUsecase(mockRepo, "secret")

	t.Run("Success Delete", func(t *testing.T) {
		mockRepo.On("Delete", uint(1)).Return(nil).Once()

		err := mockUsecase.Delete(1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
