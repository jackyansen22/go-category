package service

import (
	"context"
	"errors"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]model.Product, error)
	Search(ctx context.Context, name string, active *bool) ([]model.Product, error)
	GetByID(ctx context.Context, id int) (*model.Product, error)
	Create(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product) error
	Delete(ctx context.Context, id int) error
}

func (s *productService) Search(
	ctx context.Context,
	name string,
	active *bool,
) ([]model.Product, error) {
	return s.repo.FindByFilter(ctx, name, active)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.repo.FindAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id int) (*model.Product, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *productService) Create(
	ctx context.Context,
	p *model.Product,
) error {

	// default active
	p.Active = true

	if p.CategoryID == 0 {
		return errors.New("category_id is required")
	}

	// ✅ VALIDASI FK DI SERVICE
	if !s.repo.CategoryExists(ctx, p.CategoryID) {
		return errors.New("category not found")
	}

	return s.repo.Create(ctx, p)
}

func (s *productService) Update(ctx context.Context, p *model.Product) error {
	return s.repo.Update(ctx, p)
}

func (s *productService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
