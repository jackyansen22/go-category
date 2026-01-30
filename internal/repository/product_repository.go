package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackyansen22/crud-category/internal/model"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByID(ctx context.Context, id int) (*model.Product, error)
	Create(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(ctx context.Context) ([]model.Product, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, nama, harga, stok FROM products ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *productRepository) FindByID(ctx context.Context, id int) (*model.Product, error) {
	var p model.Product
	err := r.db.QueryRowContext(ctx,
		`SELECT id, nama, harga, stok FROM products WHERE id=$1`, id).
		Scan(&p.ID, &p.Nama, &p.Harga, &p.Stok)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *productRepository) Create(ctx context.Context, p *model.Product) error {
	return r.db.QueryRowContext(ctx,
		`INSERT INTO products (nama, harga, stok)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		p.Nama, p.Harga, p.Stok).
		Scan(&p.ID)
}

func (r *productRepository) Update(ctx context.Context, p *model.Product) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE products SET nama=$1, harga=$2, stok=$3 WHERE id=$4`,
		p.Nama, p.Harga, p.Stok, p.ID)

	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM products WHERE id=$1`, id)

	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}
	return nil
}
