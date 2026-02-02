package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type CartService struct {
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// GetOrCreateUserCart gets or creates cart for user
func (s *CartService) GetOrCreateUserCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	return s.cartRepo.GetOrCreateCart(ctx, userID)
}

// AddToCart adds a product to user's cart
func (s *CartService) AddToCart(ctx context.Context, userID uuid.UUID, req *model.AddToCartRequest) (*model.CartItemResponse, error) {
	// Validate product exists and is available
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found")
	}

	if !product.IsActive {
		return nil, fmt.Errorf("product is not available")
	}

	if product.StockQuantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: only %d available", product.StockQuantity)
	}

	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Add item to cart
	item, err := s.cartRepo.AddItem(ctx, cart.ID, req.ProductID, req.Quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to cart: %w", err)
	}

	// Get shop info
	shop, err := s.productRepo.GetShopByProductID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shop info: %w", err)
	}

	// Build response
	response := &model.CartItemResponse{
		ID:              item.ID,
		ProductID:       product.ID,
		ProductName:     product.Name,
		ProductPrice:    product.Price,
		ProductImageURL: product.ImageURL,
		StockQuantity:   product.StockQuantity,
		IsActive:        product.IsActive,
		ShopID:          product.ShopID,
		ShopName:        shop.Name,
		Quantity:        item.Quantity,
		Subtotal:        product.Price * float64(item.Quantity),
		CreatedAt:       item.CreatedAt,
		UpdatedAt:       item.UpdatedAt,
	}

	return response, nil
}

// GetCart retrieves user's cart with all items
func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*model.CartResponse, error) {
	// Get or create cart
	cart, err := s.cartRepo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	// Get cart items with product details
	items, err := s.cartRepo.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	// Build response
	response := &model.CartResponse{
		ID:        cart.ID,
		UserID:    cart.UserID,
		Items:     items,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}

	// Calculate totals
	response.CalculateTotals()

	return response, nil
}

// UpdateCartItemQuantity updates quantity of cart item
func (s *CartService) UpdateCartItemQuantity(ctx context.Context, userID, itemID uuid.UUID, quantity int) error {
	// Get cart item
	item, err := s.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("cart item not found")
	}

	// Verify ownership
	owned, err := s.cartRepo.VerifyCartOwnership(ctx, item.CartID, userID)
	if err != nil || !owned {
		return fmt.Errorf("unauthorized access to cart item")
	}

	// Validate product stock
	product, err := s.productRepo.GetByID(ctx, item.ProductID)
	if err != nil {
		return fmt.Errorf("product not found")
	}

	if product.StockQuantity < quantity {
		return fmt.Errorf("insufficient stock: only %d available", product.StockQuantity)
	}

	// Update quantity
	if err := s.cartRepo.UpdateItemQuantity(ctx, itemID, quantity); err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	return nil
}

// RemoveCartItem removes item from cart
func (s *CartService) RemoveCartItem(ctx context.Context, userID, itemID uuid.UUID) error {
	// Get cart item
	item, err := s.cartRepo.GetCartItemByID(ctx, itemID)
	if err != nil {
		return fmt.Errorf("cart item not found")
	}

	// Verify ownership
	owned, err := s.cartRepo.VerifyCartOwnership(ctx, item.CartID, userID)
	if err != nil || !owned {
		return fmt.Errorf("unauthorized access to cart item")
	}

	// Remove item
	if err := s.cartRepo.RemoveItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to remove cart item: %w", err)
	}

	return nil
}

// ClearCart removes all items from user's cart
func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("cart not found")
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// GetCartItemCount returns total items in user's cart
func (s *CartService) GetCartItemCount(ctx context.Context, userID uuid.UUID) (int, error) {
	// Get cart
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		// Return 0 if cart doesn't exist yet
		return 0, nil
	}

	return s.cartRepo.GetCartItemCount(ctx, cart.ID)
}
