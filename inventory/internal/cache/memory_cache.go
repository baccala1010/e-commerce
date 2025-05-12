package cache

import (
	"sync"
	"time"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// MemoryCache implements the ProductCache interface using in-memory storage
type MemoryCache struct {
	products        map[uuid.UUID]*model.Product
	productLists    map[string]productListCacheItem
	mu              sync.RWMutex
	refreshTicker   *time.Ticker
	stopRefreshChan chan struct{}
}

type productListCacheItem struct {
	products []model.Product
	total    int64
	created  time.Time
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		products:        make(map[uuid.UUID]*model.Product),
		productLists:    make(map[string]productListCacheItem),
		stopRefreshChan: make(chan struct{}),
	}
}

// GetProduct retrieves a product from the cache by ID
func (c *MemoryCache) GetProduct(id uuid.UUID) (*model.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	product, exists := c.products[id]
	if !exists {
		return nil, false
	}

	// Return a copy to prevent modification of cached data
	productCopy := *product
	return &productCopy, true
}

// SetProduct adds or updates a product in the cache
func (c *MemoryCache) SetProduct(product *model.Product) {
	if product == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Store a copy to prevent modification of cached data
	productCopy := *product
	c.products[product.ID] = &productCopy
}

// DeleteProduct removes a product from the cache
func (c *MemoryCache) DeleteProduct(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.products, id)
}

// GetProductList retrieves a product list from the cache by key
func (c *MemoryCache) GetProductList(key string) ([]model.Product, int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.productLists[key]
	if !exists {
		return nil, 0, false
	}

	// Return a copy to prevent modification of cached data
	productsCopy := make([]model.Product, len(item.products))
	for i, product := range item.products {
		productsCopy[i] = product
	}

	return productsCopy, item.total, true
}

// SetProductList adds or updates a product list in the cache
func (c *MemoryCache) SetProductList(key string, products []model.Product, total int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Store a copy to prevent modification of cached data
	productsCopy := make([]model.Product, len(products))
	for i, product := range products {
		productsCopy[i] = product
	}

	c.productLists[key] = productListCacheItem{
		products: productsCopy,
		total:    total,
		created:  time.Now(),
	}
}

// Clear removes all items from the cache
func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.products = make(map[uuid.UUID]*model.Product)
	c.productLists = make(map[string]productListCacheItem)
}

// StartPeriodicRefresh starts a background goroutine to refresh the cache periodically
func (c *MemoryCache) StartPeriodicRefresh(interval time.Duration, refreshFunc func() error) {
	c.stopRefreshChan = make(chan struct{})
	c.refreshTicker = time.NewTicker(interval)

	go func() {
		// Initial refresh
		if err := refreshFunc(); err != nil {
			logrus.Errorf("Error during initial cache refresh: %v", err)
		}

		for {
			select {
			case <-c.refreshTicker.C:
				if err := refreshFunc(); err != nil {
					logrus.Errorf("Error refreshing cache: %v", err)
				}
			case <-c.stopRefreshChan:
				c.refreshTicker.Stop()
				return
			}
		}
	}()

	logrus.Infof("Started periodic cache refresh every %v", interval)
}

// StopPeriodicRefresh stops the periodic refresh goroutine
func (c *MemoryCache) StopPeriodicRefresh() {
	if c.refreshTicker != nil {
		close(c.stopRefreshChan)
		c.refreshTicker = nil
		logrus.Info("Stopped periodic cache refresh")
	}
}
