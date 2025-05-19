package usecase

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/baccala1010/e-commerce/inventory/pkg/kafka"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/google/uuid"
)

type productUseCase struct {
	productRepo   repository.ProductRepository
	categoryRepo  repository.CategoryRepository
	kafkaProducer *kafka.Producer
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository, producer *kafka.Producer) ProductUseCase {
	return &productUseCase{
		productRepo:   productRepo,
		categoryRepo:  categoryRepo,
		kafkaProducer: producer,
	}
}

func (u *productUseCase) CreateProduct(request model.CreateProductRequest) (*model.Product, error) {
	// Verify category exists
	category, err := u.categoryRepo.FindByID(request.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("error finding category: %w", err)
	}

	if category == nil {
		return nil, errors.New(model.ErrCategoryNotFound)
	}

	product := &model.Product{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		StockLevel:  request.StockLevel,
		CategoryID:  request.CategoryID,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, fmt.Errorf("error creating product: %w", err)
	}

	// Publish event to Kafka if producer is available
	if u.kafkaProducer != nil {
		event := map[string]interface{}{
			"event_type":  "product_created",
			"product_id":  product.ID.String(),
			"name":        product.Name,
			"category_id": product.CategoryID.String(),
		}

		eventJSON, err := json.Marshal(event)
		if err != nil {
			// Log error but don't fail the operation
			fmt.Printf("Failed to marshal product event: %v\n", err)
		} else {
			if err := u.kafkaProducer.PublishEvent(product.ID.String(), eventJSON); err != nil {
				// Log error but don't fail the operation
				fmt.Printf("Failed to publish product created event: %v\n", err)
			} else {
				fmt.Printf("Product created event published for product: %s\n", product.ID.String())
			}
		}
	}

	return product, nil
}

func (u *productUseCase) GetProductByID(id uuid.UUID) (*model.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding product: %w", err)
	}

	if product == nil {
		return nil, errors.New(model.ErrProductNotFound)
	}

	return product, nil
}

func (u *productUseCase) UpdateProduct(id uuid.UUID, request model.UpdateProductRequest) (*model.Product, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("error finding product: %w", err)
	}

	if product == nil {
		return nil, errors.New(model.ErrProductNotFound)
	}

	// Update only provided fields
	if request.Name != nil {
		product.Name = *request.Name
	}

	if request.Description != nil {
		product.Description = *request.Description
	}

	if request.Price != nil {
		product.Price = *request.Price
	}

	if request.StockLevel != nil {
		product.StockLevel = *request.StockLevel
	}

	if request.CategoryID != nil {
		// Verify category exists
		category, err := u.categoryRepo.FindByID(*request.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("error finding category: %w", err)
		}

		if category == nil {
			return nil, errors.New(model.ErrCategoryNotFound)
		}

		product.CategoryID = *request.CategoryID
	}

	if err := u.productRepo.Update(product); err != nil {
		return nil, fmt.Errorf("error updating product: %w", err)
	}

	return product, nil
}

func (u *productUseCase) DeleteProduct(id uuid.UUID) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding product: %w", err)
	}

	if product == nil {
		return errors.New(model.ErrProductNotFound)
	}

	if err := u.productRepo.Delete(id); err != nil {
		return fmt.Errorf("error deleting product: %w", err)
	}

	return nil
}

func (u *productUseCase) ListProducts(params repository.ListProductParams) ([]model.Product, int64, error) {
	return u.productRepo.List(params)
}
