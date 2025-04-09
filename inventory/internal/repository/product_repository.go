package repository

import (
	"errors"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id uuid.UUID) (*model.Product, error)
	Update(product *model.Product) error
	Delete(id uuid.UUID) error
	List(params ListProductParams) ([]model.Product, int64, error)
}

type ListProductParams struct {
	CategoryID *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
	Search     *string
	Page       int
	PageSize   int
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uuid.UUID) (*model.Product, error) {
	var product model.Product

	if err := r.db.Preload("Category").First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Product{}, "id = ?", id).Error
}

func (r *productRepository) List(params ListProductParams) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{}).Preload("Category")

	// Apply filters
	if params.CategoryID != nil {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	if params.MinPrice != nil {
		query = query.Where("price >= ?", params.MinPrice)
	}

	if params.MaxPrice != nil {
		query = query.Where("price <= ?", params.MaxPrice)
	}

	if params.Search != nil && *params.Search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?", "%"+*params.Search+"%", "%"+*params.Search+"%")
	}

	// Count total results
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if params.Page <= 0 {
		params.Page = 1
	}

	if params.PageSize <= 0 {
		params.PageSize = 10
	}

	offset := (params.Page - 1) * params.PageSize

	if err := query.Offset(offset).Limit(params.PageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
