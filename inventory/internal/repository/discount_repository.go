package repository

import (
	"errors"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount *model.Discount) error
	FindByID(id uuid.UUID) (*model.Discount, error)
	Update(discount *model.Discount) error
	Delete(id uuid.UUID) error
	FindAll() ([]model.Discount, error)
	FindByProductID(productID uuid.UUID) ([]model.Discount, error)
	FindAllWithProducts() ([]model.Discount, error)
}

type discountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{db: db}
}

func (r *discountRepository) Create(discount *model.Discount) error {
	return r.db.Create(discount).Error
}

func (r *discountRepository) FindByID(id uuid.UUID) (*model.Discount, error) {
	var discount model.Discount

	if err := r.db.First(&discount, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &discount, nil
}

func (r *discountRepository) Update(discount *model.Discount) error {
	return r.db.Save(discount).Error
}

func (r *discountRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Discount{}, "id = ?", id).Error
}

func (r *discountRepository) FindAll() ([]model.Discount, error) {
	var discounts []model.Discount

	if err := r.db.Find(&discounts).Error; err != nil {
		return nil, err
	}

	return discounts, nil
}

func (r *discountRepository) FindByProductID(productID uuid.UUID) ([]model.Discount, error) {
	var discounts []model.Discount

	if err := r.db.Where("? = ANY(applicable_products)", productID).
		Or("applicable_products @> ?", []byte(uuid.Nil.String())).
		Find(&discounts).Error; err != nil {
		return nil, err
	}

	return discounts, nil
}

func (r *discountRepository) FindAllWithProducts() ([]model.Discount, error) {
	var discounts []model.Discount

	// Find active discounts with start_date <= current time <= end_date
	if err := r.db.Where("is_active = ? AND start_date <= NOW() AND end_date >= NOW()", true).
		Find(&discounts).Error; err != nil {
		return nil, err
	}

	return discounts, nil
}
