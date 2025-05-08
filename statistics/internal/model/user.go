package model

import "time"

// User represents a user entity in the statistics system
type User struct {
	ID               string    `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	RegistrationDate time.Time `json:"registration_date"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// UserOrderStatistics represents statistics about a user's orders
type UserOrderStatistics struct {
	UserID           string                `json:"user_id"`
	TotalOrders      int                   `json:"total_orders"`
	TotalSpent       float64               `json:"total_spent"`
	AverageOrderValue float64              `json:"average_order_value"`
	OrderTimeDistribution []OrderTimeOfDay `json:"order_time_distribution"`
	FirstOrderAt     time.Time             `json:"first_order_at"`
	LastOrderAt      time.Time             `json:"last_order_at"`
}

// OrderTimeOfDay represents the distribution of orders by hour
type OrderTimeOfDay struct {
	Hour       string `json:"hour"`
	OrderCount int    `json:"order_count"`
}

// UserStatistics represents general user statistics
type UserStatistics struct {
	TotalRegisteredUsers int `json:"total_registered_users"`
	// Add more general user statistics as needed
}