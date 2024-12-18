package service

import (
	"context"

	"forum/internal/entity"
	"forum/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.Category
}

func newCategoryService(categoryRepo repository.Category) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (r *CategoryService) GetAllCategorys(ctx context.Context) ([]entity.Category, int, error) {
	return r.categoryRepo.GetAllCategories(ctx)
}

