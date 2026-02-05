package model

import "time"

// =====================================================
// Transaction (header)
// table: transactions
// =====================================================
type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details"`
}

// =====================================================
// Transaction Detail (items)
// table: transaction_details
// =====================================================
type TransactionDetail struct {
	ID            int    `json:"id"`
	TransactionID int    `json:"transaction_id"`
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name,omitempty"` // derived (JOIN products)
	Quantity      int    `json:"quantity"`
	Subtotal      int    `json:"subtotal"`
}

// =====================================================
// Checkout Request DTO
// (NOT a database table)
// =====================================================
type CheckoutItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items"`
}
