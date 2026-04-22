package usecase_test

import (
	"testing"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/azharf99/portofolio-api/domain/mocks"
	"github.com/azharf99/portofolio-api/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPortfolioUsecase_Fetch(t *testing.T) {
	mockRepo := new(mocks.PortfolioRepositoryMock)
	mockUsecase := usecase.NewPortfolioUsecase(mockRepo)

	t.Run("Success Fetch Data", func(t *testing.T) {
		mockPortfolios := []domain.Portfolio{
			{Title: "Web Keamanan", Industry: "Cybersecurity"},
		}

		// Ekspektasi: page 2, limit 10 -> offset harusnya 10, onlyPublished true
		mockRepo.On("Fetch", 10, 10, "", "", "", true).Return(mockPortfolios, int64(1), nil).Once()

		res, total, err := mockUsecase.Fetch(2, 10, "", "", "", true)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, int64(1), total)
		assert.Equal(t, "Web Keamanan", res[0].Title)

		mockRepo.AssertExpectations(t)
	})
}

func TestPortfolioUsecase_Store(t *testing.T) {
	mockRepo := new(mocks.PortfolioRepositoryMock)
	mockUsecase := usecase.NewPortfolioUsecase(mockRepo)

	t.Run("Success Store Data", func(t *testing.T) {
		mockPortfolio := &domain.Portfolio{Title: "API Baru"}

		mockRepo.On("Store", mock.AnythingOfType("*domain.Portfolio")).Return(nil).Once()

		err := mockUsecase.Store(mockPortfolio)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
