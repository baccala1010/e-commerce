package model

import "time"

// Product represents a product entity in the statistics system
type Product struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	CategoryID string    `json:"category_id"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}