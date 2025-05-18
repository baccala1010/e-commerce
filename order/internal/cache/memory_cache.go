package cache

import (
	"sync"
	"time"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type MemoryCache struct {
	orders          map[uuid.UUID]*model.Order
	orderLists      map[string]orderListCacheItem
	mu              sync.RWMutex
	refreshTicker   *time.Ticker
	stopRefreshChan chan struct{}
}

type orderListCacheItem struct {
	orders  []model.Order
	total   int64
	created time.Time
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		orders:          make(map[uuid.UUID]*model.Order),
		orderLists:      make(map[string]orderListCacheItem),
		stopRefreshChan: make(chan struct{}),
	}
}

func (c *MemoryCache) GetOrder(id uuid.UUID) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[id]
	if !exists {
		return nil, false
	}
	orderCopy := *order
	return &orderCopy, true
}

func (c *MemoryCache) SetOrder(order *model.Order) {
	if order == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	orderCopy := *order
	c.orders[order.ID] = &orderCopy
}

func (c *MemoryCache) DeleteOrder(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.orders, id)
}

func (c *MemoryCache) GetOrderList(key string) ([]model.Order, int64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.orderLists[key]
	if !exists {
		return nil, 0, false
	}
	ordersCopy := make([]model.Order, len(item.orders))
	for i, order := range item.orders {
		ordersCopy[i] = order
	}
	return ordersCopy, item.total, true
}

func (c *MemoryCache) SetOrderList(key string, orders []model.Order, total int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ordersCopy := make([]model.Order, len(orders))
	for i, order := range orders {
		ordersCopy[i] = order
	}
	c.orderLists[key] = orderListCacheItem{
		orders:  ordersCopy,
		total:   total,
		created: time.Now(),
	}
}

func (c *MemoryCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders = make(map[uuid.UUID]*model.Order)
	c.orderLists = make(map[string]orderListCacheItem)
}

func (c *MemoryCache) StartPeriodicRefresh(interval time.Duration, refreshFunc func() error) {
	c.stopRefreshChan = make(chan struct{})
	c.refreshTicker = time.NewTicker(interval)
	go func() {
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

func (c *MemoryCache) StopPeriodicRefresh() {
	if c.refreshTicker != nil {
		close(c.stopRefreshChan)
		c.refreshTicker = nil
		logrus.Info("Stopped periodic cache refresh")
	}
}
