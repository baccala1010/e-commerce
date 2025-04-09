package repository

import (
	"errors"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindByID(id uuid.UUID) (*model.Category, error)
	Update(category *model.Category) error
	Delete(id uuid.UUID) error
	FindAll() ([]model.Category, error)
	HasProducts(id uuid.UUID) (bool, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*model.Category, error) {
	var category model.Category

	if err := r.db.First(&category, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Category{}, "id = ?", id).Error
}

func (r *categoryRepository) HasProducts(id uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.Model(&model.Product{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}
