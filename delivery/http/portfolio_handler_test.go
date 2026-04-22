package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	// Kita tambahkan alias "apiHttp" di depannya
	apiHttp "github.com/azharf99/portofolio-api/delivery/http"
	"github.com/azharf99/portofolio-api/domain"
	"github.com/azharf99/portofolio-api/domain/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPortfolioHandler_Fetch(t *testing.T) {
	// Set Gin ke mode test agar log tidak berisik
	gin.SetMode(gin.TestMode)

	mockUsecase := new(mocks.PortfolioUsecaseMock)

	// Setup Router Gin Palsu
	r := gin.Default()

	// Gunakan alias apiHttp yang baru kita buat
	pH := apiHttp.NewPortfolioHandlerInstance(mockUsecase)
	r.GET("/portfolios", pH.Fetch)

	t.Run("Success Get Portfolios", func(t *testing.T) {
		mockData := []domain.Portfolio{
			{Title: "Sistem Keamanan Bank"},
		}

		// Harapannya: dipanggil dengan default page=1, limit=10, dan onlyPublished=true
		mockUsecase.On("Fetch", 1, 10, "", "", "", true).Return(mockData, int64(1), nil).Once()

		// Buat Request HTTP Palsu (menggunakan net/http bawaan)
		req, _ := http.NewRequest("GET", "/portfolios", nil)
		w := httptest.NewRecorder()

		// Eksekusi request ke router
		r.ServeHTTP(w, req)

		// Verifikasi
		assert.Equal(t, 200, w.Code)

		// Cek apakah response JSON benar
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(1), response["total"])

		mockUsecase.AssertExpectations(t)
	})
}
