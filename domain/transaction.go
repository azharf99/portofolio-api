package domain

import "time"

type Transaction struct {
	OrderID           string    `json:"order_id" gorm:"type:varchar(100);primaryKey"`
	ServiceID         uint      `json:"service_id"`
	Service           Service   `json:"service" gorm:"foreignKey:ServiceID"`
	CustomerName      string    `json:"customer_name" gorm:"type:varchar(100);not null"`
	CustomerEmail     string    `json:"customer_email" gorm:"type:varchar(100);not null"`
	CustomerPhone     string    `json:"customer_phone" gorm:"type:varchar(20)"`
	GrossAmount       int64     `json:"gross_amount"`
	TransactionStatus string    `json:"transaction_status" gorm:"type:varchar(50);default:'pending'"`
	PaymentType       string    `json:"payment_type" gorm:"type:varchar(50)"`
	PaymentURL        string    `json:"payment_url" gorm:"type:text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TransactionUsecase interface {
	Checkout(serviceID uint, name, email, phone string) (*Transaction, error)
	HandleWebhook(orderID string, status string, paymentType string) error
	FetchByEmail(email string) ([]Transaction, error)
}

type TransactionRepository interface {
	Store(tx *Transaction) error
	Update(orderID string, status string, paymentType string) error
	GetByID(orderID string) (*Transaction, error)
	FetchByEmail(email string) ([]Transaction, error)
}
