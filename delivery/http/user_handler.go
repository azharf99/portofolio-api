package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/azharf99/portofolio-api/middleware"
	i18n_pkg "github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandlerInstance(us domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: us}
}

func NewUserHandler(r *gin.RouterGroup, us domain.UserUsecase) {
	handler := &UserHandler{usecase: us}
	r.POST("/login", handler.Login)
	r.PUT("/users/:id", handler.Update)    // BARU
	r.DELETE("/users/:id", handler.Delete) // BARU
}

func (h *UserHandler) Login(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	token, err := h.usecase.Login(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": i18n_pkg.T(localizer, err.Error())})
		return
	}

	// KEAMANAN: Token disimpan sebagai cookie httpOnly, TIDAK dikembalikan di body JSON.
	// Ini mencegah frontend menaruhnya di localStorage, tempat token bisa dicuri lewat XSS.
	middleware.SetAuthCookie(c, token)
	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "success")})
}

// Logout menghapus cookie auth di sisi client.
func (h *UserHandler) Logout(c *gin.Context) {
	middleware.ClearAuthCookie(c)
	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// UpdateUserInput adalah allow-list eksplisit field yang boleh diubah client.
// KEAMANAN: Jangan bind request body langsung ke domain.User — itu membuka celah
// mass-assignment (client bisa menyisipkan id/created_at/dsb di body request).
type UpdateUserInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Update(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	user := domain.User{Username: input.Username, Password: input.Password}
	if err := h.usecase.Update(uint(id), &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")}) // or user_not_found
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "user_updated")})
}

func (h *UserHandler) Delete(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, err.Error())})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "user_deleted")})
}
