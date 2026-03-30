package http

import (
	"net/http"
	"strconv"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/gin-gonic/gin"
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
	search := c.Query("search")
	industry := c.Query("industry")
	pType := c.Query("type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	portfolios, total, err := h.usecase.Fetch(page, limit, search, industry, pType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var portfolio domain.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Store(&portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": portfolio, "message": "Portofolio berhasil ditambahkan"})
}

func (h *PortfolioHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var portfolio domain.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.Update(uint(id), &portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portofolio berhasil diperbarui"})
}

func (h *PortfolioHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	if err := h.usecase.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portofolio berhasil dihapus"})
}
