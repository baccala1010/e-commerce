package backoffice

import (
	"github.com/baccala1010/e-commerce/order/internal/model"
	"github.com/baccala1010/e-commerce/order/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Helper functions to convert between model and proto
func convertOrderToProto(order *model.Order) *pb.Order {
	return &pb.Order{
		Id:              order.ID.String(),
		UserId:          order.UserID.String(),
		Status:          convertModelOrderStatusToProto(order.Status),
		TotalAmount:     order.TotalAmount,
		ShippingName:    order.ShippingName,
		ShippingEmail:   order.ShippingEmail,
		ShippingPhone:   order.ShippingPhone,
		ShippingAddress: order.ShippingAddr,
		Payment:         convertPaymentToProto(&order.Payment),
		CreatedAt:       timestamppb.New(order.CreatedAt),
		UpdatedAt:       timestamppb.New(order.UpdatedAt),
	}
}

func convertPaymentToProto(payment *model.Payment) *pb.Payment {
	return &pb.Payment{
		Id:            payment.ID.String(),
		OrderId:       payment.OrderID.String(),
		Amount:        payment.Amount,
		Method:        convertModelPaymentMethodToProto(payment.Method),
		Status:        convertModelPaymentStatusToProto(payment.Status),
		TransactionId: payment.TransactionID,
		PaymentDate:   timestamppb.New(payment.PaymentDate),
		CreatedAt:     timestamppb.New(payment.CreatedAt),
		UpdatedAt:     timestamppb.New(payment.UpdatedAt),
	}
}

func convertModelOrderStatusToProto(status model.OrderStatus) pb.OrderStatus {
	switch status {
	case model.OrderStatusPending:
		return pb.OrderStatus_ORDER_STATUS_PENDING
	case model.OrderStatusPaid:
		return pb.OrderStatus_ORDER_STATUS_PAID
	case model.OrderStatusShipped:
		return pb.OrderStatus_ORDER_STATUS_SHIPPED
	case model.OrderStatusDelivered:
		return pb.OrderStatus_ORDER_STATUS_DELIVERED
	case model.OrderStatusCancelled:
		return pb.OrderStatus_ORDER_STATUS_CANCELLED
	default:
		return pb.OrderStatus_ORDER_STATUS_UNSPECIFIED
	}
}

func convertProtoOrderStatusToModel(status pb.OrderStatus) model.OrderStatus {
	switch status {
	case pb.OrderStatus_ORDER_STATUS_PENDING:
		return model.OrderStatusPending
	case pb.OrderStatus_ORDER_STATUS_PAID:
		return model.OrderStatusPaid
	case pb.OrderStatus_ORDER_STATUS_SHIPPED:
		return model.OrderStatusShipped
	case pb.OrderStatus_ORDER_STATUS_DELIVERED:
		return model.OrderStatusDelivered
	case pb.OrderStatus_ORDER_STATUS_CANCELLED:
		return model.OrderStatusCancelled
	default:
		return model.OrderStatusPending
	}
}

func convertModelPaymentStatusToProto(status model.PaymentStatus) pb.PaymentStatus {
	switch status {
	case model.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case model.PaymentStatusSuccess:
		return pb.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case model.PaymentStatusFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	case model.PaymentStatusRefunded:
		return pb.PaymentStatus_PAYMENT_STATUS_REFUNDED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func convertProtoPaymentStatusToModel(status pb.PaymentStatus) model.PaymentStatus {
	switch status {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return model.PaymentStatusPending
	case pb.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		return model.PaymentStatusSuccess
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return model.PaymentStatusFailed
	case pb.PaymentStatus_PAYMENT_STATUS_REFUNDED:
		return model.PaymentStatusRefunded
	default:
		return model.PaymentStatusPending
	}
}

func convertModelPaymentMethodToProto(method model.PaymentMethod) pb.PaymentMethod {
	switch method {
	case model.PaymentMethodCreditCard:
		return pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodDebitCard:
		return pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD
	case model.PaymentMethodPaypal:
		return pb.PaymentMethod_PAYMENT_METHOD_PAYPAL
	case model.PaymentMethodBankWire:
		return pb.PaymentMethod_PAYMENT_METHOD_BANK_WIRE
	default:
		return pb.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func convertProtoPaymentMethodToModel(method pb.PaymentMethod) model.PaymentMethod {
	switch method {
	case pb.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case pb.PaymentMethod_PAYMENT_METHOD_DEBIT_CARD:
		return model.PaymentMethodDebitCard
	case pb.PaymentMethod_PAYMENT_METHOD_PAYPAL:
		return model.PaymentMethodPaypal
	case pb.PaymentMethod_PAYMENT_METHOD_BANK_WIRE:
		return model.PaymentMethodBankWire
	default:
		return model.PaymentMethodCreditCard
	}
}

// Review conversion functions
func convertReviewToProto(review *model.Review) *pb.Review {
	return &pb.Review{
		Id:          review.ID.String(),
		OrderId:     review.OrderID.String(),
		UserId:      review.UserID.String(),
		Rating:      convertModelRatingToProto(review.Rating),
		Description: review.Description,
		CreateAt:    timestamppb.New(review.CreatedAt),
	}
}

func convertModelRatingToProto(rating model.Rating) pb.Rating {
	switch rating {
	case model.RatingOne:
		return pb.Rating_RATING_ONE
	case model.RatingTwo:
		return pb.Rating_RATING_TWO
	case model.RatingThree:
		return pb.Rating_RATING_THREE
	case model.RatingFour:
		return pb.Rating_RATING_FOUR
	case model.RatingFive:
		return pb.Rating_RATING_FIVE
	default:
		return pb.Rating_RATING_UNSPECIFIED
	}
}

func convertProtoRatingToModel(rating pb.Rating) model.Rating {
	switch rating {
	case pb.Rating_RATING_ONE:
		return model.RatingOne
	case pb.Rating_RATING_TWO:
		return model.RatingTwo
	case pb.Rating_RATING_THREE:
		return model.RatingThree
	case pb.Rating_RATING_FOUR:
		return model.RatingFour
	case pb.Rating_RATING_FIVE:
		return model.RatingFive
	default:
		return model.RatingOne
	}
}
