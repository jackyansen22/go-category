package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/jackyansen22/crud-category/internal/model"
)

type ProductRepository interface {
	FindAll(ctx context.Context) ([]model.Product, error)
	FindByFilter(ctx context.Context, name string, active *bool) ([]model.Product, error)
	FindByID(ctx context.Context, id int) (*model.Product, error)
	Create(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id int) error
	CategoryExists(ctx context.Context, categoryID int) bool
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

// =====================================================
// GET ALL PRODUCTS
// =====================================================
func (r *productRepository) FindAll(
	ctx context.Context,
) ([]model.Product, error) {

	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id,
			nama,
			harga,
			stok,
			active,
			category_id
		FROM products
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Nama,
			&p.Harga,
			&p.Stok,
			&p.Active,
			&p.CategoryID, // ✅ WAJIB
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// =====================================================
// GET PRODUCTS WITH FILTER (?name=&active=)
// =====================================================
func (r *productRepository) FindByFilter(
	ctx context.Context,
	name string,
	active *bool,
) ([]model.Product, error) {

	query := `
		SELECT
			id,
			nama,
			harga,
			stok,
			active,
			category_id
		FROM products
		WHERE 1=1
	`
	args := []any{}
	argPos := 1

	if name != "" {
		query += " AND LOWER(nama) LIKE LOWER($" + strconv.Itoa(argPos) + ")"
		args = append(args, "%"+name+"%")
		argPos++
	}

	if active != nil {
		query += " AND active = $" + strconv.Itoa(argPos)
		args = append(args, *active)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Nama,
			&p.Harga,
			&p.Stok,
			&p.Active,
			&p.CategoryID, // ✅ WAJIB
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

// =====================================================
// GET PRODUCT BY ID
// =====================================================
func (r *productRepository) FindByID(
	ctx context.Context,
	id int,
) (*model.Product, error) {

	var p model.Product
	err := r.db.QueryRowContext(ctx, `
		SELECT
			p.id,
			p.nama,
			p.harga,
			p.stok,
			p.active,
			p.category_id,
			c.name
		FROM products p
		JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
	`, id).Scan(
		&p.ID,
		&p.Nama,
		&p.Harga,
		&p.Stok,
		&p.Active,
		&p.CategoryID,
		&p.CategoryName,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// =====================================================
// CREATE PRODUCT
// =====================================================
func (r *productRepository) Create(
	ctx context.Context,
	p *model.Product,
) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO products
			(nama, harga, stok, active, category_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`,
		p.Nama,
		p.Harga,
		p.Stok,
		p.Active,
		p.CategoryID,
	).Scan(&p.ID)
}

// =====================================================
// UPDATE PRODUCT
// =====================================================
func (r *productRepository) Update(ctx context.Context, p *model.Product) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE products
		SET nama = $1,
		    harga = $2,
		    stok = $3,
		    active = $4
		WHERE id = $5
	`,
		p.Nama,
		p.Harga,
		p.Stok,
		p.Active,
		p.ID,
	)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

// =====================================================
// DELETE PRODUCT
// =====================================================
func (r *productRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM products
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (r *productRepository) CategoryExists(
	ctx context.Context,
	categoryID int,
) bool {

	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM categories WHERE id = $1
		)
	`, categoryID).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}
