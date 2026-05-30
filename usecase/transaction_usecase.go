package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/azharf99/portofolio-api/domain"
	"github.com/google/uuid"
)

type transactionUsecase struct {
	txRepo      domain.TransactionRepository
	serviceRepo domain.ServiceRepository
}

func NewTransactionUsecase(txRepo domain.TransactionRepository, serviceRepo domain.ServiceRepository) domain.TransactionUsecase {
	return &transactionUsecase{
		txRepo:      txRepo,
		serviceRepo: serviceRepo,
	}
}

type GatewayChargeRequest struct {
	OrderID         string           `json:"order_id"`
	GrossAmount    int64            `json:"gross_amount"`
	CustomerDetails CustomerDetails  `json:"customer_details"`
	ItemDetails    []ItemDetail     `json:"item_details"`
}

type CustomerDetails struct {
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type ItemDetail struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

type GatewayChargeResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
	OrderID     string `json:"order_id"`
	Error       string `json:"error"`
}

func (u *transactionUsecase) Checkout(serviceID uint, name, email, phone string) (*domain.Transaction, error) {
	// 1. Fetch service details
	service, err := u.serviceRepo.GetByID(serviceID)
	if err != nil {
		return nil, fmt.Errorf("service not found: %w", err)
	}

	// Determine price (promo price takes priority if > 0)
	price := service.OriginalPrice
	if service.PromoPrice > 0 {
		price = service.PromoPrice
	}

	// 2. Generate unique order ID
	orderID := fmt.Sprintf("TRX-%s", uuid.New().String())

	// 3. Prepare payload for Centralized Payment Gateway
	gatewayURL := os.Getenv("PAYMENT_GATEWAY_URL")
	gatewayAPIKey := os.Getenv("PAYMENT_GATEWAY_API_KEY")

	if gatewayURL == "" {
		return nil, errors.New("PAYMENT_GATEWAY_URL environment variable is not set")
	}
	if gatewayAPIKey == "" {
		return nil, errors.New("PAYMENT_GATEWAY_API_KEY environment variable is not set")
	}

	payload := GatewayChargeRequest{
		OrderID:     orderID,
		GrossAmount: price,
		CustomerDetails: CustomerDetails{
			FirstName: name,
			Email:     email,
			Phone:     phone,
		},
		ItemDetails: []ItemDetail{
			{
				ID:       fmt.Sprintf("srv-%d", service.ID),
				Name:     service.Title,
				Price:    price,
				Quantity: 1,
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gateway payload: %w", err)
	}

	// 4. Send request to Centralized Payment Gateway
	reqURL := fmt.Sprintf("%s/api/v1/charge", gatewayURL)
	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", gatewayAPIKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("payment gateway connection error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&errResp)
		return nil, fmt.Errorf("payment gateway returned status %d: %v", resp.StatusCode, errResp["error"])
	}

	var gatewayResp GatewayChargeResponse
	if err := json.NewDecoder(resp.Body).Decode(&gatewayResp); err != nil {
		return nil, fmt.Errorf("failed to decode payment gateway response: %w", err)
	}

	// 5. Store transaction locally
	tx := &domain.Transaction{
		OrderID:           orderID,
		ServiceID:         serviceID,
		CustomerName:      name,
		CustomerEmail:     email,
		CustomerPhone:     phone,
		GrossAmount:       price,
		TransactionStatus: "pending",
		PaymentURL:        gatewayResp.RedirectURL,
	}

	if err := u.txRepo.Store(tx); err != nil {
		return nil, fmt.Errorf("failed to save transaction locally: %w", err)
	}

	// Make sure we load the service relationship back into the struct
	tx.Service = *service

	return tx, nil
}

func (u *transactionUsecase) HandleWebhook(orderID string, status string, paymentType string) error {
	// 1. Verify transaction exists
	_, err := u.txRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	// 2. Update status in local database
	return u.txRepo.Update(orderID, status, paymentType)
}

func (u *transactionUsecase) FetchByEmail(email string) ([]domain.Transaction, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return u.txRepo.FetchByEmail(email)
}
