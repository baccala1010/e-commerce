package model

import (
	"time"

	"github.com/google/uuid"
)

// UserOrderStatistic stores statistics about user orders
type UserOrderStatistic struct {
	UserID     uuid.UUID `gorm:"primaryKey;type:uuid;not null" json:"user_id"`
	OrderCount int       `gorm:"not null;default:0" json:"order_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
