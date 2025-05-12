package cache

import (
	"time"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/google/uuid"
)

// ProductCache defines the interface for caching product operations
type ProductCache interface {
	// Product operations
	GetProduct(id uuid.UUID) (*model.Product, bool)
	SetProduct(product *model.Product)
	DeleteProduct(id uuid.UUID)

	// Product list operations
	GetProductList(key string) ([]model.Product, int64, bool)
	SetProductList(key string, products []model.Product, total int64)

	// Cache management
	Clear()

	// Start background refresh
	StartPeriodicRefresh(interval time.Duration, refreshFunc func() error)
	StopPeriodicRefresh()
}
