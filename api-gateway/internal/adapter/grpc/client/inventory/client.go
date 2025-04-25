package inventory

import (
	"context"
	"fmt"

	"github.com/baccala1010/e-commerce/api-gateway/internal/config"
	"github.com/baccala1010/e-commerce/api-gateway/pkg/grpcconn"
	"github.com/baccala1010/e-commerce/inventory/pkg/pb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

// Client is a gRPC client for the inventory service
type Client struct {
	conn   *grpc.ClientConn
	client pb.InventoryServiceClient
}

// NewClient creates a new inventory service client
func NewClient(ctx context.Context, connManager *grpcconn.ConnectionManager, cfg *config.Config) (*Client, error) {
	conn, err := connManager.GetConnection(ctx, "inventory", cfg.Services.Inventory)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory service connection: %w", err)
	}

	client := pb.NewInventoryServiceClient(conn)
	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

// CreateProduct creates a new product
func (c *Client) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	logrus.Infof("Calling inventory service CreateProduct: %+v", req)
	return c.client.CreateProduct(ctx, req)
}

// GetProductByID gets a product by ID
func (c *Client) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	logrus.Infof("Calling inventory service GetProductByID: %+v", req)
	return c.client.GetProductByID(ctx, req)
}

// UpdateProduct updates a product
func (c *Client) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	logrus.Infof("Calling inventory service UpdateProduct: %+v", req)
	return c.client.UpdateProduct(ctx, req)
}

// DeleteProduct deletes a product
func (c *Client) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) error {
	logrus.Infof("Calling inventory service DeleteProduct: %+v", req)
	_, err := c.client.DeleteProduct(ctx, req)
	return err
}

// ListProducts lists products
func (c *Client) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	logrus.Infof("Calling inventory service ListProducts: %+v", req)
	return c.client.ListProducts(ctx, req)
}

// CreateCategory creates a new category
func (c *Client) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	logrus.Infof("Calling inventory service CreateCategory: %+v", req)
	return c.client.CreateCategory(ctx, req)
}

// GetCategoryByID gets a category by ID
func (c *Client) GetCategoryByID(ctx context.Context, req *pb.GetCategoryRequest) (*pb.CategoryResponse, error) {
	logrus.Infof("Calling inventory service GetCategoryByID: %+v", req)
	return c.client.GetCategoryByID(ctx, req)
}

// UpdateCategory updates a category
func (c *Client) UpdateCategory(ctx context.Context, req *pb.UpdateCategoryRequest) (*pb.CategoryResponse, error) {
	logrus.Infof("Calling inventory service UpdateCategory: %+v", req)
	return c.client.UpdateCategory(ctx, req)
}

// DeleteCategory deletes a category
func (c *Client) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) error {
	logrus.Infof("Calling inventory service DeleteCategory: %+v", req)
	_, err := c.client.DeleteCategory(ctx, req)
	return err
}

// ListCategories lists categories
func (c *Client) ListCategories(ctx context.Context, req *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	logrus.Infof("Calling inventory service ListCategories: %+v", req)
	return c.client.ListCategories(ctx, req)
}

// CreateDiscount creates a new discount
func (c *Client) CreateDiscount(ctx context.Context, req *pb.CreateDiscountRequest) (*pb.DiscountResponse, error) {
	logrus.Infof("Calling inventory service CreateDiscount: %+v", req)
	return c.client.CreateDiscount(ctx, req)
}

// GetDiscountByID gets a discount by ID
func (c *Client) GetDiscountByID(ctx context.Context, req *pb.GetDiscountRequest) (*pb.DiscountResponse, error) {
	logrus.Infof("Calling inventory service GetDiscountByID: %+v", req)
	return c.client.GetDiscountByID(ctx, req)
}

// UpdateDiscount updates a discount
func (c *Client) UpdateDiscount(ctx context.Context, req *pb.UpdateDiscountRequest) (*pb.DiscountResponse, error) {
	logrus.Infof("Calling inventory service UpdateDiscount: %+v", req)
	return c.client.UpdateDiscount(ctx, req)
}

// DeleteDiscount deletes a discount
func (c *Client) DeleteDiscount(ctx context.Context, req *pb.DeleteDiscountRequest) error {
	logrus.Infof("Calling inventory service DeleteDiscount: %+v", req)
	_, err := c.client.DeleteDiscount(ctx, req)
	return err
}

// GetAllProductsWithPromotion gets all products with promotions
func (c *Client) GetAllProductsWithPromotion(ctx context.Context, req *pb.GetProductsWithPromotionRequest) (*pb.ListProductsResponse, error) {
	logrus.Infof("Calling inventory service GetAllProductsWithPromotion: %+v", req)
	return c.client.GetAllProductsWithPromotion(ctx, req)
}

// GetProductsByDiscountID gets products by discount ID
func (c *Client) GetProductsByDiscountID(ctx context.Context, req *pb.GetProductsByDiscountIDRequest) (*pb.ListProductsResponse, error) {
	logrus.Infof("Calling inventory service GetProductsByDiscountID: %+v", req)
	return c.client.GetProductsByDiscountID(ctx, req)
}
