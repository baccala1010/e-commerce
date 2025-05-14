package cache

import (
	"time"

	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/google/uuid"
)

type OrderCache interface {
	GetOrder(id uuid.UUID) (*model.Order, bool)
	SetOrder(order *model.Order)
	DeleteOrder(id uuid.UUID)
	GetOrderList(key string) ([]model.Order, int64, bool)
	SetOrderList(key string, orders []model.Order, total int64)
	Clear()
	StartPeriodicRefresh(interval time.Duration, refreshFunc func() error)
	StopPeriodicRefresh()
}
