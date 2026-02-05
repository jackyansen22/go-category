package model

type Product struct {
	ID           int    `json:"id"`
	Nama         string `json:"nama"`
	Harga        int    `json:"harga"`
	Stok         int    `json:"stok"`
	Active       bool   `json:"active"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name,omitempty"` // JOIN result
}
