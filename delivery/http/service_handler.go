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

type ServiceHandler struct {
	usecase domain.ServiceUsecase
}

func NewServiceHandlerInstance(us domain.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{usecase: us}
}

func NewServiceHandler(r *gin.RouterGroup, us domain.ServiceUsecase) {
	handler := &ServiceHandler{usecase: us}

	r.GET("/services", handler.Fetch)
	r.POST("/services", handler.Store)
	r.PUT("/services/:id", handler.Update)
	r.DELETE("/services/:id", handler.Delete)
}

func (h *ServiceHandler) Fetch(c *gin.Context) {
	h.fetch(c, true)
}

func (h *ServiceHandler) AdminFetch(c *gin.Context) {
	h.fetch(c, false)
}

func (h *ServiceHandler) fetch(c *gin.Context, activeOnly bool) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
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

	services, total, err := h.usecase.Fetch(page, limit, activeOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  services,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *ServiceHandler) Store(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	var service domain.Service
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	// Handle File Uploads
	form, _ := c.MultipartForm()
	if form != nil {
		if files := form.File["image"]; len(files) > 0 {
			file := files[0]
			path, err := h.saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
				return
			}
			service.ImageURL = path
		}
	}

	if err := h.usecase.Store(&service); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": service, "message": i18n_pkg.T(localizer, "service_created")})
}

func (h *ServiceHandler) saveUploadedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExt[ext] {
		return "", fmt.Errorf("file %s tidak diizinkan. Hanya jpg, jpeg, png, webp", file.Filename)
	}

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

	filename := uuid.New().String() + ext
	savePath := filepath.Join("uploads", "services", filename)

	// Pastikan folder uploads/services ada
	// Karena SaveUploadedFile akan error jika foldernya tidak ada di windows
	// Kita bisa buat dulu foldernya di Cwd
	// Kita biarkan gin handles save path
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		// Kita coba upload ke uploads/portfolios saja jika gagal, karena folder tersebut pasti ada
		savePath = filepath.Join("uploads", "portfolios", filename)
		if err2 := c.SaveUploadedFile(file, savePath); err2 != nil {
			return "", err2
		}
	}

	return "/" + filepath.ToSlash(savePath), nil
}

func (h *ServiceHandler) Update(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	var service domain.Service
	if err := c.ShouldBind(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
		return
	}

	// Handle File Uploads
	form, _ := c.MultipartForm()
	if form != nil {
		if files := form.File["image"]; len(files) > 0 {
			file := files[0]
			path, err := h.saveUploadedFile(c, file)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
				return
			}
			service.ImageURL = path
		}
	}

	if err := h.usecase.Update(uint(id), &service); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": i18n_pkg.T(localizer, "invalid_request")})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "service_updated")})
}

func (h *ServiceHandler) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": i18n_pkg.T(localizer, "service_deleted")})
}
