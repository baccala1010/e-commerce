package repository

import (
	"fmt"
	"log"

	"github.com/baccala1010/e-commerce/order/internal/cache"
	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
)

type CachedOrderRepository struct {
	repo  OrderRepository
	cache cache.OrderCache
}

func NewCachedOrderRepository(repo OrderRepository, cache cache.OrderCache) OrderRepository {
	return &CachedOrderRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedOrderRepository) Create(order *model.Order) error {
	err := r.repo.Create(order)
	if err != nil {
		return err
	}
	r.cache.SetOrder(order)
	return nil
}

func (r *CachedOrderRepository) FindByID(id uuid.UUID) (*model.Order, error) {
	if order, found := r.cache.GetOrder(id); found {
		log.Printf("[CACHE HIT] Order ID: %s", id)
		return order, nil
	}
	log.Printf("[CACHE MISS] Order ID: %s", id)
	order, err := r.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order != nil {
		r.cache.SetOrder(order)
	}
	return order, nil
}

// Remove ListOrderParams and List method, wrap FindByUserID for caching
func (r *CachedOrderRepository) FindByUserID(userID uuid.UUID, page, pageSize int) ([]model.Order, int64, error) {
	cacheKey := fmt.Sprintf("user:%s:page:%d:size:%d", userID.String(), page, pageSize)
	if orders, total, found := r.cache.GetOrderList(cacheKey); found {
		log.Printf("[CACHE HIT] Order List Key: %s", cacheKey)
		return orders, total, nil
	}
	log.Printf("[CACHE MISS] Order List Key: %s", cacheKey)
	orders, total, err := r.repo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	r.cache.SetOrderList(cacheKey, orders, total)
	return orders, total, nil
}

func (r *CachedOrderRepository) RefreshCache() error {
	r.cache.Clear()
	// Example: load all orders (for demo, first 1000 orders)
	// In production, use paging or a more scalable approach
	orders, err := r.repo.FindAll()
	if err != nil {
		return err
	}
	for i := range orders {
		r.cache.SetOrder(&orders[i])
	}
	log.Printf("[CACHE INIT] Loaded %d orders into cache", len(orders))
	return nil
}

// Add missing Update method to satisfy OrderRepository interface
func (r *CachedOrderRepository) Update(order *model.Order) error {
	err := r.repo.Update(order)
	if err != nil {
		return err
	}
	r.cache.SetOrder(order)
	return nil
}

func (r *CachedOrderRepository) FindAll() ([]model.Order, error) {
	return r.repo.FindAll()
}
