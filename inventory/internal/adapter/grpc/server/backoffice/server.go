package backoffice

import (
	"context"

	"github.com/baccala1010/e-commerce/inventory/internal/model"
	"github.com/baccala1010/e-commerce/inventory/internal/repository"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Product methods
func (s *Server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
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

	product, err := s.productUseCase.CreateProduct(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return &pb.ProductResponse{
		Product: convertProductToProto(product),
	}, nil
}

func (s *Server) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	product, err := s.productUseCase.GetProductByID(productID)
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

func (s *Server) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
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

	product, err := s.productUseCase.UpdateProduct(productID, updateReq)
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

func (s *Server) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	productID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product ID: %v", err)
	}

	if err := s.productUseCase.DeleteProduct(productID); err != nil {
		if err.Error() == model.ErrProductNotFound {
			return nil, status.Errorf(codes.NotFound, "product not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
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

	products, total, err := s.productUseCase.ListProducts(params)
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
func (s *Server) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	createReq := model.CreateCategoryRequest{
		Name:        req.Name,
		Description: req.Description,
	}

	category, err := s.categoryUseCase.CreateCategory(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create category: %v", err)
	}

	return &pb.CategoryResponse{
		Category: convertCategoryToProto(category),
	}, nil
}

func (s *Server) GetCategoryByID(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	categoryID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	category, err := s.categoryUseCase.GetCategoryByID(categoryID)
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

func (s *Server) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
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

	category, err := s.categoryUseCase.UpdateCategory(categoryID, updateReq)
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

func (s *Server) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*emptypb.Empty, error) {
	categoryID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid category ID: %v", err)
	}

	if err := s.categoryUseCase.DeleteCategory(categoryID); err != nil {
		if err.Error() == model.ErrCategoryNotFound {
			return nil, status.Errorf(codes.NotFound, "category not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete category: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	categories, err := s.categoryUseCase.ListCategories()
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

// Discount methods
func (s *Server) CreateDiscount(ctx context.Context, req *pb.CreateDiscountRequest) (*pb.DiscountResponse, error) {
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

	discount, err := s.discountUseCase.CreateDiscount(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create discount: %v", err)
	}

	return &pb.DiscountResponse{
		Discount: convertDiscountToProto(discount),
	}, nil
}

func (s *Server) GetDiscountByID(ctx context.Context, req *pb.GetDiscountRequest) (*pb.DiscountResponse, error) {
	discountID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	discount, err := s.discountUseCase.GetDiscountByID(discountID)
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

func (s *Server) UpdateDiscount(ctx context.Context, req *pb.UpdateDiscountRequest) (*pb.DiscountResponse, error) {
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

	discount, err := s.discountUseCase.UpdateDiscount(discountID, updateReq)
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

func (s *Server) DeleteDiscount(ctx context.Context, req *pb.DeleteDiscountRequest) (*emptypb.Empty, error) {
	discountID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	err = s.discountUseCase.DeleteDiscount(discountID)
	if err != nil {
		if err.Error() == model.ErrDiscountNotFound {
			return nil, status.Errorf(codes.NotFound, "discount not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete discount: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *Server) GetAllProductsWithPromotion(ctx context.Context, req *pb.GetProductsWithPromotionRequest) (*pb.ListProductsResponse, error) {
	productsWithPromotions, err := s.discountUseCase.GetAllProductsWithPromotion()
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

func (s *Server) GetProductsByDiscountID(ctx context.Context, req *pb.GetProductsByDiscountIDRequest) (*pb.ListProductsResponse, error) {
	discountID, err := uuid.Parse(req.DiscountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid discount ID: %v", err)
	}

	products, err := s.discountUseCase.GetProductsByDiscountID(discountID)
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
