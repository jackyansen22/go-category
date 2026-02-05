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

	FindAll(ctx context.Context) ([]model.Transaction, error)
	FindByID(ctx context.Context, id int) (*model.Transaction, error)
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

func (r *transactionRepository) FindAll(
	ctx context.Context,
) ([]model.Transaction, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, total_amount, created_at
		FROM transactions
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.TotalAmount,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionRepository) FindByID(
	ctx context.Context,
	id int,
) (*model.Transaction, error) {

	var t model.Transaction
	err := r.db.QueryRowContext(ctx, `
		SELECT id, total_amount, created_at
		FROM transactions
		WHERE id = $1
	`, id).Scan(&t.ID, &t.TotalAmount, &t.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("transaction not found")
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			td.id,
			td.transaction_id,
			td.product_id,
			p.nama,
			td.quantity,
			td.subtotal
		FROM transaction_details td
		JOIN products p ON p.id = td.product_id
		WHERE td.transaction_id = $1
		ORDER BY td.id
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var d model.TransactionDetail
		if err := rows.Scan(
			&d.ID,
			&d.TransactionID,
			&d.ProductID,
			&d.ProductName,
			&d.Quantity,
			&d.Subtotal,
		); err != nil {
			return nil, err
		}
		t.Details = append(t.Details, d)
	}

	return &t, nil
}
