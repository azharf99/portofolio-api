package repository

import (
	"github.com/azharf99/portofolio-api/domain"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) Store(tx *domain.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) Update(orderID string, status string, paymentType string) error {
	return r.db.Model(&domain.Transaction{}).
		Where("order_id = ?", orderID).
		Updates(map[string]interface{}{
			"transaction_status": status,
			"payment_type":       paymentType,
		}).Error
}

func (r *transactionRepository) GetByID(orderID string) (*domain.Transaction, error) {
	var tx domain.Transaction
	err := r.db.Preload("Service").Where("order_id = ?", orderID).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FetchByEmail(email string) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	err := r.db.Preload("Service").Where("customer_email = ?", email).Order("created_at DESC").Find(&txs).Error
	return txs, err
}
