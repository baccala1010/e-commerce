package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusSuccess  PaymentStatus = "success"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodDebitCard  PaymentMethod = "debit_card"
	PaymentMethodPaypal     PaymentMethod = "paypal"
	PaymentMethodBankWire   PaymentMethod = "bank_wire"
)

type Payment struct {
	ID            uuid.UUID      `json:"id" gorm:"type:uuid;primary_key"`
	OrderID       uuid.UUID      `json:"order_id" gorm:"type:uuid;not null"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	Method        PaymentMethod  `json:"method" gorm:"type:varchar(20);not null"`
	Status        PaymentStatus  `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	TransactionID string         `json:"transaction_id" gorm:"type:varchar(255)"`
	PaymentDate   time.Time      `json:"payment_date"`
	CreatedAt     time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"not null;default:now()"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (p *Payment) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	if p.Status == "" {
		p.Status = PaymentStatusPending
	}

	return nil
}

// PaymentDTO is used for order creation requests
type PaymentDTO struct {
	Method PaymentMethod `json:"method" binding:"required,oneof=credit_card debit_card paypal bank_wire"`
}

// ProcessPaymentRequest represents request for processing a payment
type ProcessPaymentRequest struct {
	OrderID uuid.UUID     `json:"order_id" binding:"required"`
	Method  PaymentMethod `json:"method" binding:"required,oneof=credit_card debit_card paypal bank_wire"`
}

// UpdatePaymentStatusRequest represents request for updating payment status
type UpdatePaymentStatusRequest struct {
	Status        PaymentStatus `json:"status" binding:"required,oneof=pending success failed refunded"`
	TransactionID string        `json:"transaction_id"`
}
