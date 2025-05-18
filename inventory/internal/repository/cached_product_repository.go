package repository

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/baccala1010/e-commerce/inventory/internal/cache"
	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
)

// CachedProductRepository implements ProductRepository with caching
type CachedProductRepository struct {
	repo  ProductRepository
	cache cache.ProductCache
}

// NewCachedProductRepository creates a new cached product repository
func NewCachedProductRepository(repo ProductRepository, cache cache.ProductCache) ProductRepository {
	return &CachedProductRepository{
		repo:  repo,
		cache: cache,
	}
}

// Create creates a new product and adds it to the cache
func (r *CachedProductRepository) Create(product *model.Product) error {
	// Create the product in the database
	err := r.repo.Create(product)
	if err != nil {
		return err
	}

	// Add the product to the cache
	r.cache.SetProduct(product)

	return nil
}

// FindByID retrieves a product by ID, using the cache if available
func (r *CachedProductRepository) FindByID(id uuid.UUID) (*model.Product, error) {
	// Try to get the product from the cache
	if product, found := r.cache.GetProduct(id); found {
		return product, nil
	}

	// If not in cache, get from the database
	product, err := r.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// If product was found, add it to the cache
	if product != nil {
		r.cache.SetProduct(product)
	}

	return product, nil
}

// Update updates a product and updates it in the cache
func (r *CachedProductRepository) Update(product *model.Product) error {
	// Update the product in the database
	err := r.repo.Update(product)
	if err != nil {
		return err
	}

	// Update the product in the cache
	r.cache.SetProduct(product)

	return nil
}

// Delete deletes a product and removes it from the cache
func (r *CachedProductRepository) Delete(id uuid.UUID) error {
	// Delete the product from the database
	err := r.repo.Delete(id)
	if err != nil {
		return err
	}

	// Remove the product from the cache
	r.cache.DeleteProduct(id)

	return nil
}

// List retrieves a list of products, using the cache if available
func (r *CachedProductRepository) List(params ListProductParams) ([]model.Product, int64, error) {
	// Generate a cache key based on the params
	cacheKey, err := generateCacheKey(params)
	if err != nil {
		// If we can't generate a cache key, just use the repository directly
		return r.repo.List(params)
	}

	// Try to get the product list from the cache
	if products, total, found := r.cache.GetProductList(cacheKey); found {
		return products, total, nil
	}

	// If not in cache, get from the database
	products, total, err := r.repo.List(params)
	if err != nil {
		return nil, 0, err
	}

	// Add the product list to the cache
	r.cache.SetProductList(cacheKey, products, total)

	return products, total, nil
}

// RefreshCache refreshes the cache with all products
func (r *CachedProductRepository) RefreshCache() error {
	// Clear the cache
	r.cache.Clear()

	// Get all products
	products, total, err := r.repo.List(ListProductParams{
		Page:     1,
		PageSize: 1000, // Use a large page size to get all products
	})
	if err != nil {
		return fmt.Errorf("error refreshing cache: %w", err)
	}

	// Cache individual products
	for i := range products {
		r.cache.SetProduct(&products[i])
	}

	// Cache the full product list
	r.cache.SetProductList("all", products, total)

	return nil
}

// generateCacheKey creates a unique key for caching product lists based on the params
func generateCacheKey(params ListProductParams) (string, error) {
	// Convert params to JSON
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	// Create an MD5 hash of the JSON
	hash := md5.Sum(paramsJSON)
	return fmt.Sprintf("products:%s", hex.EncodeToString(hash[:])), nil
}
