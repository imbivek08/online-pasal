package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: productRepo,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, vendorID uuid.UUID, req *model.CreateProductRequest) (*model.Product, error) {
	product := &model.Product{
		ID:            uuid.New(),
		VendorID:      vendorID,
		Name:          req.Name,
		Description:   req.Description,
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		Category:      req.Category,
		ImageURL:      req.ImageURL,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllProducts retrieves all products
func (s *ProductService) GetAllProducts(ctx context.Context, filters map[string]interface{}) ([]*model.Product, error) {
	return s.repo.GetAll(ctx, filters)
}

// GetVendorProducts retrieves all products for a vendor
func (s *ProductService) GetVendorProducts(ctx context.Context, vendorID uuid.UUID) ([]*model.Product, error) {
	return s.repo.GetByVendorID(ctx, vendorID)
}

// UpdateProduct updates a product
func (s *ProductService) UpdateProduct(ctx context.Context, productID uuid.UUID, vendorID uuid.UUID, req *model.UpdateProductRequest) (*model.Product, error) {
	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Check if vendor owns the product
	if product.VendorID != vendorID {
		return nil, err // You could create a custom error for unauthorized access
	}

	// Update fields if provided
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.StockQuantity != nil {
		product.StockQuantity = *req.StockQuantity
	}
	if req.Category != nil {
		product.Category = req.Category
	}
	if req.ImageURL != nil {
		product.ImageURL = req.ImageURL
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product
func (s *ProductService) DeleteProduct(ctx context.Context, productID uuid.UUID, vendorID uuid.UUID) error {
	// Get existing product
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	// Check if vendor owns the product
	if product.VendorID != vendorID {
		return err // You could create a custom error for unauthorized access
	}

	return s.repo.Delete(ctx, productID)
}
