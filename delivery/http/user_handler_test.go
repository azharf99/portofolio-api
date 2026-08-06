package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	apiHttp "github.com/azharf99/portofolio-api/delivery/http"
	"github.com/azharf99/portofolio-api/domain"
	"github.com/azharf99/portofolio-api/domain/mocks"
	"github.com/azharf99/portofolio-api/middleware"
	"github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Login(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.Init()

	mockUsecase := new(mocks.UserUsecaseMock)
	r := gin.Default()
	r.Use(middleware.I18nMiddleware())
	uH := apiHttp.NewUserHandlerInstance(mockUsecase)
	r.POST("/login", uH.Login)

	t.Run("Success Login", func(t *testing.T) {
		mockToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

		// Harapan: usecase.Login mengembalikan token
		mockUsecase.On("Login", "azharfa", "admin123").Return(mockToken, nil).Once()

		// Buat body JSON untuk request
		requestBody, _ := json.Marshal(map[string]string{
			"username": "azharfa",
			"password": "admin123",
		})

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		// Token TIDAK boleh muncul di body JSON (harus lewat cookie httpOnly saja).
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.NotContains(t, response, "token")

		// Token harus ada di cookie httpOnly bernama middleware.AuthCookieName.
		var authCookie *http.Cookie
		for _, c := range w.Result().Cookies() {
			if c.Name == middleware.AuthCookieName {
				authCookie = c
			}
		}
		if assert.NotNil(t, authCookie, "expected httpOnly auth cookie to be set") {
			assert.Equal(t, mockToken, authCookie.Value)
			assert.True(t, authCookie.HttpOnly)
		}

		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.Init()

	mockUsecase := new(mocks.UserUsecaseMock)
	r := gin.Default()
	r.Use(middleware.I18nMiddleware())
	uH := apiHttp.NewUserHandlerInstance(mockUsecase)
	r.PUT("/users/:id", uH.Update)

	t.Run("Success Update", func(t *testing.T) {
		mockUser := domain.User{Username: "azharfa_update"}

		mockUsecase.On("Update", uint(1), &mockUser).Return(nil).Once()

		requestBody, _ := json.Marshal(mockUser)
		req, _ := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestUserHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	i18n.Init()

	mockUsecase := new(mocks.UserUsecaseMock)
	r := gin.Default()
	r.Use(middleware.I18nMiddleware())
	uH := apiHttp.NewUserHandlerInstance(mockUsecase)
	r.DELETE("/users/:id", uH.Delete)

	t.Run("Success Delete", func(t *testing.T) {
		mockUsecase.On("Delete", uint(1)).Return(nil).Once()

		req, _ := http.NewRequest("DELETE", "/users/1", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}
