package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackyansen22/crud-category/internal/model"
)

type TransactionRepository interface {
	CreateTransaction(
		ctx context.Context,
		items []model.CheckoutItem,
	) (*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// =====================================================
// CREATE TRANSACTION (CHECKOUT)
// - atomic
// - row locking (FOR UPDATE)
// - stock validation
// =====================================================
func (r *transactionRepository) CreateTransaction(
	ctx context.Context,
	items []model.CheckoutItem,
) (*model.Transaction, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]model.TransactionDetail, 0)

	for _, item := range items {
		var (
			productName  string
			productPrice int
			stock        int
		)

		// 🔒 lock product row
		err := tx.QueryRowContext(ctx, `
			SELECT nama, harga, stok
			FROM products
			WHERE id = $1
			FOR UPDATE
		`, item.ProductID).Scan(&productName, &productPrice, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf(
				"stock not enough for product %d (available %d)",
				item.ProductID, stock,
			)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		// update stock
		_, err = tx.ExecContext(ctx, `
			UPDATE products
			SET stok = stok - $1
			WHERE id = $2
		`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, model.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	if totalAmount <= 0 {
		return nil, errors.New("total amount must be greater than zero")
	}

	var transactionID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO transactions (total_amount)
		VALUES ($1)
		RETURNING id
	`, totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID

		_, err = tx.ExecContext(ctx, `
			INSERT INTO transaction_details
				(transaction_id, product_id, quantity, subtotal)
			VALUES ($1, $2, $3, $4)
		`,
			transactionID,
			details[i].ProductID,
			details[i].Quantity,
			details[i].Subtotal,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
