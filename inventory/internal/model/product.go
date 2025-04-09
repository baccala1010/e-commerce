package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	StockLevel  int            `json:"stock_level" gorm:"not null"`
	CategoryID  uuid.UUID      `json:"category_id" gorm:"type:uuid;not null"`
	Category    Category       `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

// CreateProductRequest represents the request body for creating a new product
type CreateProductRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	StockLevel  int       `json:"stock_level" binding:"required,gte=0"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
}

// UpdateProductRequest represents the request body for updating a product
type UpdateProductRequest struct {
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
	Price       *float64   `json:"price" binding:"omitempty,gt=0"`
	StockLevel  *int       `json:"stock_level" binding:"omitempty,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id"`
}
