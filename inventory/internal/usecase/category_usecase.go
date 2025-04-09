package usecase

import (
	"errors"
	"fmt"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/google/uuid"
)

type categoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryUseCase creates a new category use case
func NewCategoryUseCase(categoryRepo repository.CategoryRepository) CategoryUseCase {
	return &categoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUseCase) CreateCategory(request model.CreateCategoryRequest) (*model.Category, error) {
	category := &model.Category{
		Name:        request.Name,
		Description: request.Description,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, fmt.Errorf("error creating category: %w", err)
	}

	return category, nil
}

func (u *categoryUseCase) GetCategoryByID(id uuid.UUID) (*model.Category, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding category: %w", err)
	}

	if category == nil {
		return nil, errors.New(model.ErrCategoryNotFound)
	}

	return category, nil
}

func (u *categoryUseCase) UpdateCategory(id uuid.UUID, request model.UpdateCategoryRequest) (*model.Category, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding category: %w", err)
	}

	if category == nil {
		return nil, errors.New(model.ErrCategoryNotFound)
	}

	if request.Name != nil {
		category.Name = *request.Name
	}

	if request.Description != nil {
		category.Description = *request.Description
	}

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, fmt.Errorf("error updating category: %w", err)
	}

	return category, nil
}

func (u *categoryUseCase) DeleteCategory(id uuid.UUID) error {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding category: %w", err)
	}

	if category == nil {
		return errors.New(model.ErrCategoryNotFound)
	}

	// Check if there are products using this category
	hasProducts, err := u.categoryRepo.HasProducts(id)
	if err != nil {
		return fmt.Errorf("error checking for products: %w", err)
	}

	if hasProducts {
		return errors.New("cannot delete category with associated products")
	}

	if err := u.categoryRepo.Delete(id); err != nil {
		return fmt.Errorf("error deleting category: %w", err)
	}

	return nil
}

func (u *categoryUseCase) ListCategories() ([]model.Category, error) {
	return u.categoryRepo.FindAll()
}
