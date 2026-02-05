package service

import (
	"context"
	"errors"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/repository"
)

type TransactionService interface {
	Checkout(ctx context.Context, items []model.CheckoutItem) (*model.Transaction, error)

	GetAll(ctx context.Context) ([]model.Transaction, error)
	GetByID(ctx context.Context, id int) (*model.Transaction, error)
}

type transactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// =====================================================
// CHECKOUT
// - atomic transaction
// - calculate subtotal & total
// =====================================================
func (s *transactionService) Checkout(
	ctx context.Context,
	items []model.CheckoutItem,
) (*model.Transaction, error) {

	if len(items) == 0 {
		return nil, errors.New("checkout items cannot be empty")
	}

	// Business orchestration delegated to repository (sql.Tx)
	return s.repo.CreateTransaction(ctx, items)
}

func (s *transactionService) GetAll(
	ctx context.Context,
) ([]model.Transaction, error) {
	return s.repo.FindAll(ctx)
}

func (s *transactionService) GetByID(
	ctx context.Context,
	id int,
) (*model.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}
