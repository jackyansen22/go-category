package service

import (
	"context"

	"github.com/jackyansen22/crud-category/internal/model"
	"github.com/jackyansen22/crud-category/internal/repository"
)

type CategoryService interface {
	GetAll(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id int) (*model.Category, error)
	Create(ctx context.Context, c *model.Category) error
	Update(ctx context.Context, c *model.Category) error
	Delete(ctx context.Context, id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAll(ctx context.Context) ([]model.Category, error) {
	return s.repo.FindAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *categoryService) Create(ctx context.Context, c *model.Category) error {
	return s.repo.Create(ctx, c)
}

func (s *categoryService) Update(ctx context.Context, c *model.Category) error {
	return s.repo.Update(ctx, c)
}

func (s *categoryService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
