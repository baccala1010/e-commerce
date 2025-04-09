package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type Order struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	UserID        uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Status        OrderStatus    `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TotalAmount   float64        `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	ShippingName  string         `json:"shipping_name" gorm:"type:varchar(255);not null"`
	ShippingEmail string         `json:"shipping_email" gorm:"type:varchar(255);not null"`
	ShippingPhone string         `json:"shipping_phone" gorm:"type:varchar(20);not null"`
	ShippingAddr  string         `json:"shipping_address" gorm:"type:text;not null"`
	Payment       Payment        `json:"payment" gorm:"foreignKey:OrderID"`
	CreatedAt     time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}

	if o.Status == "" {
		o.Status = OrderStatusPending
	}

	return nil
}

// CreateOrderRequest represents the request body for creating a new order
type CreateOrderRequest struct {
	UserID        uuid.UUID  `json:"user_id" binding:"required"`
	TotalAmount   float64    `json:"total_amount" binding:"required,gt=0"`
	Payment       PaymentDTO `json:"payment" binding:"required"`
	ShippingName  string     `json:"shipping_name" binding:"required"`
	ShippingEmail string     `json:"shipping_email" binding:"required,email"`
	ShippingPhone string     `json:"shipping_phone" binding:"required"`
	ShippingAddr  string     `json:"shipping_address" binding:"required"`
}

// UpdateOrderStatusRequest represents the request body for updating an order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=pending paid shipped delivered cancelled"`
}

// ProductInfo is used to store product details from the inventory service
type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	StockLevel  int       `json:"stock_level"`
	CategoryID  uuid.UUID `json:"category_id"`
	Category    struct {
		ID          uuid.UUID `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
	} `json:"category"`
}
