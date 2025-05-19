package model

import "github.com/google/uuid"

type UserOrderStatistic struct {
	UserID     uuid.UUID `db:"user_id"`
	OrderCount int       `db:"order_count"`
}
