package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UUIDArray is a custom type for handling arrays of UUIDs in JSONB format
type UUIDArray []uuid.UUID

// Value implements the driver.Valuer interface
func (u UUIDArray) Value() (driver.Value, error) {
	if len(u) == 0 {
		return "[]", nil
	}
	return json.Marshal(u)
}

// Scan implements the sql.Scanner interface
func (u *UUIDArray) Scan(value interface{}) error {
	if value == nil {
		*u = UUIDArray{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, u)
}

type Discount struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	Name               string         `json:"name" gorm:"type:varchar(255);not null"`
	Description        string         `json:"description" gorm:"type:text"`
	DiscountPercentage float64        `json:"discount_percentage" gorm:"type:decimal(5,2);not null"`
	ApplicableProducts UUIDArray      `json:"applicable_products" gorm:"type:jsonb"`
	StartDate          time.Time      `json:"start_date" gorm:"not null"`
	EndDate            time.Time      `json:"end_date" gorm:"not null"`
	IsActive           bool           `json:"is_active" gorm:"default:true"`
	CreatedAt          time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt          time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt          gorm.DeletedAt `json:"-" gorm:"index"`
}

func (d *Discount) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}

type CreateDiscountRequest struct {
	Name               string      `json:"name" binding:"required"`
	Description        string      `json:"description"`
	DiscountPercentage float64     `json:"discount_percentage" binding:"required,gt=0,lte=100"`
	ApplicableProducts []uuid.UUID `json:"applicable_products"`
	StartDate          time.Time   `json:"start_date" binding:"required"`
	EndDate            time.Time   `json:"end_date" binding:"required"`
}

type UpdateDiscountRequest struct {
	Name               *string     `json:"name"`
	Description        *string     `json:"description"`
	DiscountPercentage *float64    `json:"discount_percentage" binding:"omitempty,gt=0,lte=100"`
	ApplicableProducts []uuid.UUID `json:"applicable_products"`
	StartDate          *time.Time  `json:"start_date"`
	EndDate            *time.Time  `json:"end_date"`
	IsActive           *bool       `json:"is_active"`
}

type DiscountResponse struct {
	ID                 uuid.UUID   `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	DiscountPercentage float64     `json:"discount_percentage"`
	ApplicableProducts []uuid.UUID `json:"applicable_products"`
	StartDate          time.Time   `json:"start_date"`
	EndDate            time.Time   `json:"end_date"`
	IsActive           bool        `json:"is_active"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}
