package usecase

import (
	"errors"
	"fmt"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/google/uuid"
)

type DiscountUseCase interface {
	CreateDiscount(request model.CreateDiscountRequest) (*model.Discount, error)
	GetDiscountByID(id uuid.UUID) (*model.Discount, error)
	UpdateDiscount(id uuid.UUID, request model.UpdateDiscountRequest) (*model.Discount, error)
	DeleteDiscount(id uuid.UUID) error
	ListDiscounts() ([]model.Discount, error)
	GetAllProductsWithPromotion() ([]ProductWithPromotions, error)
	GetProductsByDiscountID(discountID uuid.UUID) ([]model.Product, error)
}

type ProductWithPromotions struct {
	Product   model.Product    `json:"product"`
	Discounts []model.Discount `json:"discounts"`
}

type discountUseCase struct {
	discountRepo repository.DiscountRepository
	productRepo  repository.ProductRepository
}

// NewDiscountUseCase creates a new discount use case
func NewDiscountUseCase(discountRepo repository.DiscountRepository, productRepo repository.ProductRepository) DiscountUseCase {
	return &discountUseCase{
		discountRepo: discountRepo,
		productRepo:  productRepo,
	}
}

func (u *discountUseCase) CreateDiscount(request model.CreateDiscountRequest) (*model.Discount, error) {
	discount := &model.Discount{
		Name:               request.Name,
		Description:        request.Description,
		DiscountPercentage: request.DiscountPercentage,
		ApplicableProducts: model.UUIDArray(request.ApplicableProducts),
		StartDate:          request.StartDate,
		EndDate:            request.EndDate,
		IsActive:           true,
	}

	if err := u.discountRepo.Create(discount); err != nil {
		return nil, fmt.Errorf("error creating discount: %w", err)
	}

	return discount, nil
}

func (u *discountUseCase) GetDiscountByID(id uuid.UUID) (*model.Discount, error) {
	discount, err := u.discountRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding discount: %w", err)
	}

	if discount == nil {
		return nil, errors.New(model.ErrDiscountNotFound)
	}

	return discount, nil
}

func (u *discountUseCase) UpdateDiscount(id uuid.UUID, request model.UpdateDiscountRequest) (*model.Discount, error) {
	discount, err := u.discountRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding discount: %w", err)
	}

	if discount == nil {
		return nil, errors.New(model.ErrDiscountNotFound)
	}

	// Update only provided fields
	if request.Name != nil {
		discount.Name = *request.Name
	}

	if request.Description != nil {
		discount.Description = *request.Description
	}

	if request.DiscountPercentage != nil {
		discount.DiscountPercentage = *request.DiscountPercentage
	}

	if request.ApplicableProducts != nil {
		discount.ApplicableProducts = model.UUIDArray(request.ApplicableProducts)
	}

	if request.StartDate != nil {
		discount.StartDate = *request.StartDate
	}

	if request.EndDate != nil {
		discount.EndDate = *request.EndDate
	}

	if request.IsActive != nil {
		discount.IsActive = *request.IsActive
	}

	if err := u.discountRepo.Update(discount); err != nil {
		return nil, fmt.Errorf("error updating discount: %w", err)
	}

	return discount, nil
}

func (u *discountUseCase) DeleteDiscount(id uuid.UUID) error {
	discount, err := u.discountRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding discount: %w", err)
	}

	if discount == nil {
		return errors.New(model.ErrDiscountNotFound)
	}

	if err := u.discountRepo.Delete(id); err != nil {
		return fmt.Errorf("error deleting discount: %w", err)
	}

	return nil
}

func (u *discountUseCase) ListDiscounts() ([]model.Discount, error) {
	return u.discountRepo.FindAll()
}

func (u *discountUseCase) GetAllProductsWithPromotion() ([]ProductWithPromotions, error) {
	// Get all active discounts
	discounts, err := u.discountRepo.FindAllWithProducts()
	if err != nil {
		return nil, fmt.Errorf("error finding discounts: %w", err)
	}

	// Get all products
	products, _, err := u.productRepo.List(repository.ListProductParams{})
	if err != nil {
		return nil, fmt.Errorf("error finding products: %w", err)
	}

	// Map products with their applicable discounts
	result := make([]ProductWithPromotions, 0)

	for _, product := range products {
		productWithDiscounts := ProductWithPromotions{
			Product:   product,
			Discounts: []model.Discount{},
		}

		// Check which discounts apply to this product
		for _, discount := range discounts {
			// Check if the discount applies to this product
			isApplicable := false

			// If there are no applicable products, it applies to all
			if len(discount.ApplicableProducts) == 0 {
				isApplicable = true
			} else {
				// Check if product ID is in the applicable products list
				for _, applicableProductID := range discount.ApplicableProducts {
					if product.ID == applicableProductID {
						isApplicable = true
						break
					}
				}
			}

			if isApplicable {
				productWithDiscounts.Discounts = append(productWithDiscounts.Discounts, discount)
			}
		}

		// Only add products that have at least one discount
		if len(productWithDiscounts.Discounts) > 0 {
			result = append(result, productWithDiscounts)
		}
	}

	return result, nil
}

func (u *discountUseCase) GetProductsByDiscountID(discountID uuid.UUID) ([]model.Product, error) {
	// First, verify the discount exists
	discount, err := u.discountRepo.FindByID(discountID)
	if err != nil {
		return nil, fmt.Errorf("error finding discount: %w", err)
	}

	if discount == nil {
		return nil, errors.New(model.ErrDiscountNotFound)
	}

	// If there are no applicable products, it means the discount applies to all products
	if len(discount.ApplicableProducts) == 0 {
		products, _, err := u.productRepo.List(repository.ListProductParams{})
		if err != nil {
			return nil, fmt.Errorf("error finding products: %w", err)
		}
		return products, nil
	}

	// Otherwise, get only the products from the applicable products list
	result := []model.Product{}

	for _, productID := range discount.ApplicableProducts {
		product, err := u.productRepo.FindByID(productID)
		if err != nil {
			continue // Skip errors for individual products
		}

		if product != nil {
			result = append(result, *product)
		}
	}

	return result, nil
}
