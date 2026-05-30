package http

import (
	"net/http"
	"os"

	"github.com/azharf99/portofolio-api/domain"
	i18n_pkg "github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type TransactionHandler struct {
	usecase domain.TransactionUsecase
}

func NewTransactionHandlerInstance(us domain.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{usecase: us}
}

type CheckoutInput struct {
	ServiceID uint   `json:"service_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
}

type WebhookInput struct {
	OrderID           string `json:"order_id" binding:"required"`
	TransactionStatus string `json:"transaction_status" binding:"required"`
	PaymentType       string `json:"payment_type"`
}

func (h *TransactionHandler) Checkout(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	var input CheckoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": i18n_pkg.T(localizer, "invalid_request") + ": " + err.Error()})
		return
	}

	tx, err := h.usecase.Checkout(input.ServiceID, input.Name, input.Email, input.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    tx,
		"message": i18n_pkg.T(localizer, "success"),
	})
}

func (h *TransactionHandler) Webhook(c *gin.Context) {
	// Verify request is from Centralized Payment Gateway by checking X-API-Key
	apiKey := c.GetHeader("X-API-Key")
	expectedKey := os.Getenv("PAYMENT_GATEWAY_API_KEY")

	if expectedKey == "" || apiKey != expectedKey {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized webhook source"})
		return
	}

	var input WebhookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	err := h.usecase.HandleWebhook(input.OrderID, input.TransactionStatus, input.PaymentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle webhook: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully"})
}

func (h *TransactionHandler) FetchHistory(c *gin.Context) {
	localizer := c.MustGet("localizer").(*i18n.Localizer)
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email query param is required"})
		return
	}

	txs, err := h.usecase.FetchByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": i18n_pkg.T(localizer, "internal_error")})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    txs,
		"message": i18n_pkg.T(localizer, "success"),
	})
}
