package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	repo *repositories.TransactionRepository
	// kalau butuh productRepo untuk validasi lebih lanjut, bisa ditambah di sini
	// productRepo *repositories.ProductRepository
}

func NewTransactionService(repo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	if len(items) == 0 {
		return nil, errors.New("tidak ada item dalam checkout")
	}

	// Optional: tambah validasi lain kalau perlu
	// misal cek duplikat product_id, quantity > 0, dll

	return s.repo.CreateTransaction(items)
}