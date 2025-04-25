package backoffice

import (
	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Helper functions to convert between model and proto
func convertProductToProto(product *model.Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		StockLevel:  int32(product.StockLevel),
		CategoryId:  product.CategoryID.String(),
		Category:    convertCategoryToProto(&product.Category),
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}

func convertCategoryToProto(category *model.Category) *pb.Category {
	return &pb.Category{
		Id:          category.ID.String(),
		Name:        category.Name,
		Description: category.Description,
		CreatedAt:   timestamppb.New(category.CreatedAt),
		UpdatedAt:   timestamppb.New(category.UpdatedAt),
	}
}

// convertDiscountToProto converts a model.Discount to a pb.Discount
func convertDiscountToProto(discount *model.Discount) *pb.Discount {
	// Convert applicable products from UUIDs to strings
	applicableProducts := make([]string, 0, len(discount.ApplicableProducts))
	for _, productID := range discount.ApplicableProducts {
		applicableProducts = append(applicableProducts, productID.String())
	}

	return &pb.Discount{
		Id:                 discount.ID.String(),
		Name:               discount.Name,
		Description:        discount.Description,
		DiscountPercentage: discount.DiscountPercentage,
		ApplicableProducts: applicableProducts,
		StartDate:          timestamppb.New(discount.StartDate),
		EndDate:            timestamppb.New(discount.EndDate),
		IsActive:           discount.IsActive,
		CreatedAt:          timestamppb.New(discount.CreatedAt),
		UpdatedAt:          timestamppb.New(discount.UpdatedAt),
	}
}
