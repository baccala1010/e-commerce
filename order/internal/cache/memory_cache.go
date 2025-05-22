package cache

import (
	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
	"sync"
	"time"
)

type InMemoryOrderCache struct {
	mu         sync.RWMutex
	items      map[uuid.UUID]*model.Order
	listCache  map[string][]model.Order
	totalCache map[string]int64
	stopCh     chan struct{}
}

func NewInMemoryOrderCache() *InMemoryOrderCache {
	return &InMemoryOrderCache{
		items:      make(map[uuid.UUID]*model.Order),
		listCache:  make(map[string][]model.Order),
		totalCache: make(map[string]int64),
		stopCh:     make(chan struct{}),
	}
}

func (c *InMemoryOrderCache) GetOrder(id uuid.UUID) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.items[id]
	return order, ok
}

func (c *InMemoryOrderCache) SetOrder(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[order.ID] = order
}

func (c *InMemoryOrderCache) DeleteOrder(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

func (c *InMemoryOrderCache) GetOrderList(key string) ([]model.Order, int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	list, ok := c.listCache[key]
	total := c.totalCache[key]
	return list, total, ok
}

func (c *InMemoryOrderCache) SetOrderList(key string, orders []model.Order, total int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.listCache[key] = orders
	c.totalCache[key] = total
}

func (c *InMemoryOrderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[uuid.UUID]*model.Order)
	c.listCache = make(map[string][]model.Order)
	c.totalCache = make(map[string]int64)
}

func (c *InMemoryOrderCache) StartPeriodicRefresh(interval time.Duration, refreshFunc func() error) {
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

func (c *InMemoryOrderCache) StopPeriodicRefresh() {
	close(c.stopCh)
}
