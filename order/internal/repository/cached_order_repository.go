package repository

import (
	"fmt"

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
		return order, nil
	}
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
		return orders, total, nil
	}
	orders, total, err := r.repo.FindByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	r.cache.SetOrderList(cacheKey, orders, total)
	return orders, total, nil
}

func (r *CachedOrderRepository) RefreshCache() error {
	// Example: refresh cache for first 100 users (in production, use a better approach)
	// This is a placeholder; you may want to load all orders or recent orders
	return nil // Implement as needed
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
