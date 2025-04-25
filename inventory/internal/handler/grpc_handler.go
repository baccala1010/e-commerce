package handler

import (
	"context"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/baccala1010/e-commerce/inventory/internal/usecase"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	pb.UnimplementedInventoryServiceServer
	productUseCase  usecase.ProductUseCase
	categoryUseCase usecase.CategoryUseCase
	discountUseCase usecase.DiscountUseCase
}

func NewGRPCHandler(productUseCase usecase.ProductUseCase, categoryUseCase usecase.CategoryUseCase, discountUseCase usecase.DiscountUseCase) *GRPCHandler {
	return &GRPCHandler{
		productUseCase:  productUseCase,
		categoryUseCase: categoryUseCase,
		discountUseCase: discountUseCase,
	}
}

// Product methods
func (h *GRPCHandler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	categoryID, err := uuid.Parse(req.CategoryId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	createReq := model.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		StockLevel:  int(req.StockLevel),
		CategoryID:  categoryID,
	}

	product, err := h.productUseCase.CreateProduct(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.ProductResponse{
		Product: convertProductToProto(product),
	}, nil
}

func (h *GRPCHandler) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	product, err := h.productUseCase.GetProductByID(productID)
	if err != nil {
		if err.Error() == model.ErrProductNotFound {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	return &pb.ProductResponse{
		Product: convertProductToProto(product),
	}, nil
}

func (h *GRPCHandler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	updateReq := model.UpdateProductRequest{}
	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.Description != nil {
		updateReq.Description = req.Description
	}
	if req.Price != nil {
		updateReq.Price = req.Price
	}
	if req.StockLevel != nil {
		stockLevel := int(*req.StockLevel)
		updateReq.StockLevel = &stockLevel
	}
	if req.CategoryId != nil {
		categoryID, err := uuid.Parse(*req.CategoryId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
		}
		updateReq.CategoryID = &categoryID
	}

	product, err := h.productUseCase.UpdateProduct(productID, updateReq)
	if err != nil {
		if err.Error() == model.ErrProductNotFound {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		if err.Error() == model.ErrCategoryNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}

	return &pb.ProductResponse{
		Product: convertProductToProto(product),
	}, nil
}

func (h *GRPCHandler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	if err := h.productUseCase.DeleteProduct(productID); err != nil {
		if err.Error() == model.ErrProductNotFound {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	params := repository.ListProductParams{
		Page:     int(req.Page),
		PageSize: int(req.Limit),
	}

	if req.CategoryId != "" {
		categoryID, err := uuid.Parse(req.CategoryId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
		}
		params.CategoryID = &categoryID
	}

	products, total, err := h.productUseCase.ListProducts(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list products: %v", err)
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = convertProductToProto(&product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
	}, nil
}

// Category methods
func (h *GRPCHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	createReq := model.CreateCategoryRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	category, err := h.categoryUseCase.CreateCategory(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	return &pb.CategoryResponse{
		Category: convertCategoryToProto(category),
	}, nil
}

func (h *GRPCHandler) GetCategoryByID(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	categoryID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	category, err := h.categoryUseCase.GetCategoryByID(categoryID)
	if err != nil {
		if err.Error() == model.ErrCategoryNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get category: %v", err)
	}

	return &pb.CategoryResponse{
		Category: convertCategoryToProto(category),
	}, nil
}

func (h *GRPCHandler) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	categoryID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	updateReq := model.UpdateCategoryRequest{}
	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.Description != nil {
		updateReq.Description = req.Description
	}

	category, err := h.categoryUseCase.UpdateCategory(categoryID, updateReq)
	if err != nil {
		if err.Error() == model.ErrCategoryNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update category: %v", err)
	}

	return &pb.CategoryResponse{
		Category: convertCategoryToProto(category),
	}, nil
}

func (h *GRPCHandler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	categoryID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	if err := h.categoryUseCase.DeleteCategory(categoryID); err != nil {
		if err.Error() == model.ErrCategoryNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, err := h.categoryUseCase.ListCategories()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list categories: %v", err)
	}

	protoCategories := make([]*pb.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = convertCategoryToProto(&category)
	}

	return &pb.ListCategoriesResponse{
		Categories: protoCategories,
		Total:      int32(len(categories)),
	}, nil
}

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

// Discount methods
func (h *GRPCHandler) CreateDiscount(ctx context.Context, req *pb.CreateDiscountRequest) (*pb.DiscountResponse, error) {
	// Convert applicable products from strings to UUIDs
	applicableProducts := make([]uuid.UUID, 0, len(req.ApplicableProducts))
	for _, productID := range req.ApplicableProducts {
		id, err := uuid.Parse(productID)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
		}
		applicableProducts = append(applicableProducts, id)
	}

	createReq := model.CreateDiscountRequest{
		Name:               req.Name,
		Description:        req.Description,
		DiscountPercentage: req.DiscountPercentage,
		ApplicableProducts: applicableProducts,
		StartDate:          req.StartDate.AsTime(),
		EndDate:            req.EndDate.AsTime(),
	}

	discount, err := h.discountUseCase.CreateDiscount(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create discount: %v", err)
	}

	return &pb.DiscountResponse{
		Discount: convertDiscountToProto(discount),
	}, nil
}

func (h *GRPCHandler) GetDiscountByID(ctx context.Context, req *pb.GetDiscountRequest) (*pb.DiscountResponse, error) {
	discountID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	discount, err := h.discountUseCase.GetDiscountByID(discountID)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			return nil, status.Errorf(codes.NotFound, "discount not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get discount: %v", err)
	}

	return &pb.DiscountResponse{
		Discount: convertDiscountToProto(discount),
	}, nil
}

func (h *GRPCHandler) UpdateDiscount(ctx context.Context, req *pb.UpdateDiscountRequest) (*pb.DiscountResponse, error) {
	discountID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	updateReq := model.UpdateDiscountRequest{}
	if req.Name != nil {
		updateReq.Name = req.Name
	}
	if req.Description != nil {
		updateReq.Description = req.Description
	}
	if req.DiscountPercentage != nil {
		updateReq.DiscountPercentage = req.DiscountPercentage
	}
	if req.IsActive != nil {
		updateReq.IsActive = req.IsActive
	}

	// Convert applicable products from strings to UUIDs
	if len(req.ApplicableProducts) > 0 {
		applicableProducts := make([]uuid.UUID, 0, len(req.ApplicableProducts))
		for _, productID := range req.ApplicableProducts {
			id, err := uuid.Parse(productID)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
			}
			applicableProducts = append(applicableProducts, id)
		}
		updateReq.ApplicableProducts = applicableProducts
	}

	if req.StartDate != nil {
		startDate := req.StartDate.AsTime()
		updateReq.StartDate = &startDate
	}
	if req.EndDate != nil {
		endDate := req.EndDate.AsTime()
		updateReq.EndDate = &endDate
	}

	discount, err := h.discountUseCase.UpdateDiscount(discountID, updateReq)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			return nil, status.Errorf(codes.NotFound, "discount not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update discount: %v", err)
	}

	return &pb.DiscountResponse{
		Discount: convertDiscountToProto(discount),
	}, nil
}

func (h *GRPCHandler) DeleteDiscount(ctx context.Context, req *pb.DeleteDiscountRequest) (*emptypb.Empty, error) {
	discountID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	err = h.discountUseCase.DeleteDiscount(discountID)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			return nil, status.Errorf(codes.NotFound, "discount not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete discount: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *GRPCHandler) GetAllProductsWithPromotion(ctx context.Context, req *pb.GetProductsWithPromotionRequest) (*pb.ListProductsResponse, error) {
	productsWithPromotions, err := h.discountUseCase.GetAllProductsWithPromotion()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products with promotions: %v", err)
	}

	// Extract just the products from the result
	products := make([]model.Product, 0, len(productsWithPromotions))
	for _, pwp := range productsWithPromotions {
		products = append(products, pwp.Product)
	}

	// Convert to proto format
	protoProducts := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		protoProducts = append(protoProducts, convertProductToProto(&product))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(protoProducts)),
	}, nil
}

func (h *GRPCHandler) GetProductsByDiscountID(ctx context.Context, req *pb.GetProductsByDiscountIDRequest) (*pb.ListProductsResponse, error) {
	discountID, err := uuid.Parse(req.DiscountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	products, err := h.discountUseCase.GetProductsByDiscountID(discountID)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			return nil, status.Errorf(codes.NotFound, "discount not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get products by discount ID: %v", err)
	}

	// Convert to proto format
	protoProducts := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		protoProducts = append(protoProducts, convertProductToProto(&product))
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(len(protoProducts)),
	}, nil
}

// Helper function to convert model.Discount to pb.Discount
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
