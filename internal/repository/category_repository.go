package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackyansen22/crud-category/internal/model"
)

type CategoryRepository interface {
	FindAll(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id int) (*model.Category, error)
	Create(ctx context.Context, c *model.Category) error
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, description
		FROM categories
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id int) (*model.Category, error) {
	var c model.Category

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, description
		FROM categories
		WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("category not found")
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *categoryRepository) Create(ctx context.Context, c *model.Category) error {
	return r.db.QueryRowContext(ctx, `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id
	`, c.Name, c.Description).Scan(&c.ID)
}

func (r *categoryRepository) Update(ctx context.Context, c *model.Category) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE categories
		SET name = $1, description = $2
		WHERE id = $3
	`, c.Name, c.Description, c.ID)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM categories
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("category not found")
	}

	return nil
}
