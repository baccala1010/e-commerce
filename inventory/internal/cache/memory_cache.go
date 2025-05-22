package cache

import (
	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
	"sync"
	"time"
)

type InMemoryProductCache struct {
	mu    sync.RWMutex
	items map[uuid.UUID]*model.Product
	// For list caching (optional, simple implementation)
	listCache  map[string][]model.Product
	totalCache map[string]int64
	stopCh     chan struct{}
}

func NewInMemoryProductCache() *InMemoryProductCache {
	return &InMemoryProductCache{
		items:      make(map[uuid.UUID]*model.Product),
		listCache:  make(map[string][]model.Product),
		totalCache: make(map[string]int64),
		stopCh:     make(chan struct{}),
	}
}

// --- Product operations ---
func (c *InMemoryProductCache) GetProduct(id uuid.UUID) (*model.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.items[id]
	return p, ok
}

func (c *InMemoryProductCache) SetProduct(product *model.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[product.ID] = product
}

func (c *InMemoryProductCache) DeleteProduct(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

// --- Product list operations ---
func (c *InMemoryProductCache) GetProductList(key string) ([]model.Product, int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list, ok := c.listCache[key]
	total := c.totalCache[key]
	return list, total, ok
}

func (c *InMemoryProductCache) SetProductList(key string, products []model.Product, total int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listCache[key] = products
	c.totalCache[key] = total
}

// --- Cache management ---
func (c *InMemoryProductCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[uuid.UUID]*model.Product)
	c.listCache = make(map[string][]model.Product)
	c.totalCache = make(map[string]int64)
}

// --- Periodic refresh ---
func (c *InMemoryProductCache) StartPeriodicRefresh(interval time.Duration, refreshFunc func() error) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = refreshFunc()
			case <-c.stopCh:
				return
			}
		}
	}()
}

func (c *InMemoryProductCache) StopPeriodicRefresh() {
	close(c.stopCh)
}
