package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null;unique"`
	Description string         `json:"description" gorm:"type:text"`
	Products    []Product      `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// CreateCategoryRequest represents the request body for creating a new category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
