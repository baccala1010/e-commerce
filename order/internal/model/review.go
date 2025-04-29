package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Rating string

const (
	RatingOne   Rating = "1"
	RatingTwo   Rating = "2"
	RatingThree Rating = "3"
	RatingFour  Rating = "4"
	RatingFive  Rating = "5"
)

type Review struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Rating      Rating    `json:"rating" gorm:"type:varchar(20);not null"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
}

func (r *Review) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}

	if r.Rating == "" {
		r.Rating = RatingOne
	}

	return nil
}

type CreateReviewRequest struct {
	OrderID     uuid.UUID `json:"order_id" binding:"required"`
	UserID      uuid.UUID `json:"user_id" binding:"required"`
	Rating      Rating    `json:"rating" binding:"required,oneof=1 2 3 4 5"`
	Description string    `json:"description" binding:"required"`
}

type GetReviewRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

type GetOrderReviewsRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

type GetOrderReviewsResponse struct {
	Reviews []Review `json:"reviews"`
}

type ReviewResponse struct {
	Review Review `json:"review"`
}
