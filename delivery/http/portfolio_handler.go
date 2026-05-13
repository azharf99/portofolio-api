package http

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/azharf99/portofolio-api/domain"
	i18n_pkg "github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

type PortfolioHandler struct {
	usecase domain.PortfolioUsecase
}

func NewPortfolioHandlerInstance(us domain.PortfolioUsecase) *PortfolioHandler {
	return &PortfolioHandler{usecase: us}
}

func NewPortfolioHandler(r *gin.RouterGroup, us domain.PortfolioUsecase) {
	handler := &PortfolioHandler{usecase: us}

	r.GET("/portfolios", handler.Fetch)
	r.POST("/portfolios", handler.Store)
	r.PUT("/portfolios/:id", handler.Update)
	r.DELETE("/portfolios/:id", handler.Delete)
}

func (h *PortfolioHandler) Fetch(c *gin.Context) {
	h.fetch(c, true)
}

func (h *PortfolioHandler) AdminFetch(c *gin.Context) {
	h.fetch(c, false)
}

func (h *PortfolioHandler) fetch(c *gin.Context, onlyPublished bool) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	search := c.Query("search")
	industry := c.Query("industry")
	pType := c.Query("type")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "page_numeric")})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "limit_numeric")})
		return
	}

	if page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "page_positive")})
		return
	}
	if limit <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "limit_positive")})
		return
	}

	portfolios, total, err := h.usecase.Fetch(page, limit, search, industry, pType, onlyPublished)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  portfolios,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *PortfolioHandler) Store(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	var portfolio domain.Portfolio
	// Gunakan ShouldBind agar bisa menangani JSON maupun Form Data
	if err := c.ShouldBind(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	// Handle File Uploads
	form, _ := c.MultipartForm()
	if form != nil {
		// 1. Handle Main Image (single)
		if files := form.File["image"]; len(files) > 0 {
			file := files[0]
			path, err := h.saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
				return
			}
			portfolio.ImageURL = path
		}

		// 2. Handle Gallery Images (multiple)
		if files := form.File["images"]; len(files) > 0 {
			var gallery []domain.PortfolioImage
			for _, file := range files {
				path, err := h.saveUploadedFile(c, file)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
					return
				}
				gallery = append(gallery, domain.PortfolioImage{ImageURL: path})
			}
			portfolio.Images = gallery
		}
	}

	if err := h.usecase.Store(&portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": portfolio, "message": i18n_pkg.T(localizer, "portfolio_created")})
}

func (h *PortfolioHandler) saveUploadedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	// 1. Filter: Cek Ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExt[ext] {
		return "", fmt.Errorf("file %s tidak diizinkan. Hanya jpg, jpeg, png, webp", file.Filename)
	}

	// 2. Filter: Cek MIME Type (Deteksi konten berbahaya)
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buffer := make([]byte, 512)
	_, err = f.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("file %s bukan merupakan file gambar yang valid", file.Filename)
	}

	// 3. Cegah Konflik Nama: Gunakan UUID
	filename := uuid.New().String() + ext
	savePath := filepath.Join("uploads", "portfolios", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		return "", err
	}

	// Kembalikan URL path (misal: /uploads/portfolios/...)
	return "/" + filepath.ToSlash(savePath), nil
}

func (h *PortfolioHandler) Update(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	var portfolio domain.Portfolio
	if err := c.ShouldBind(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	// Handle File Uploads
	form, _ := c.MultipartForm()
	if form != nil {
		// 1. Handle Main Image (single)
		if files := form.File["image"]; len(files) > 0 {
			file := files[0]
			path, err := h.saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
				return
			}
			portfolio.ImageURL = path
		}

		// 2. Handle Gallery Images (multiple)
		if files := form.File["images"]; len(files) > 0 {
			var gallery []domain.PortfolioImage
			for _, file := range files {
				path, err := h.saveUploadedFile(c, file)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
					return
				}
				gallery = append(gallery, domain.PortfolioImage{ImageURL: path})
			}
			portfolio.Images = gallery
		}
	}

	if err := h.usecase.Update(uint(id), &portfolio); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "portfolio_updated")})
}

func (h *PortfolioHandler) Delete(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "portfolio_deleted")})
}
