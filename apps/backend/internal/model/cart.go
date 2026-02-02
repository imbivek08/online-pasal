package model

import (
	"time"

	"github.com/google/uuid"
)

// Cart represents a shopping cart
type Cart struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	SessionID *string    `json:"session_id,omitempty" db:"session_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CartID    uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID uuid.UUID `json:"product_id" db:"product_id"`
	Quantity  int       `json:"quantity" db:"quantity"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CartItemWithProduct represents a cart item with product details
type CartItemWithProduct struct {
	ID              uuid.UUID `json:"id" db:"id"`
	CartID          uuid.UUID `json:"cart_id" db:"cart_id"`
	ProductID       uuid.UUID `json:"product_id" db:"product_id"`
	ProductName     string    `json:"product_name" db:"product_name"`
	ProductPrice    float64   `json:"product_price" db:"product_price"`
	ProductImageURL *string   `json:"product_image_url,omitempty" db:"product_image_url"`
	StockQuantity   int       `json:"stock_quantity" db:"stock_quantity"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	ShopID          uuid.UUID `json:"shop_id" db:"shop_id"`
	ShopName        string    `json:"shop_name" db:"shop_name"`
	Quantity        int       `json:"quantity" db:"quantity"`
	Subtotal        float64   `json:"subtotal"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// AddToCartRequest represents request to add item to cart
type AddToCartRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

// UpdateCartItemRequest represents request to update cart item quantity
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

// CartResponse represents the cart with all items
type CartResponse struct {
	ID        uuid.UUID             `json:"id"`
	UserID    *uuid.UUID            `json:"user_id,omitempty"`
	Items     []CartItemWithProduct `json:"items"`
	ItemCount int                   `json:"item_count"`
	Subtotal  float64               `json:"subtotal"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// CartItemResponse represents a single cart item response
type CartItemResponse struct {
	ID              uuid.UUID `json:"id"`
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	ProductPrice    float64   `json:"product_price"`
	ProductImageURL *string   `json:"product_image_url,omitempty"`
	StockQuantity   int       `json:"stock_quantity"`
	IsActive        bool      `json:"is_active"`
	ShopID          uuid.UUID `json:"shop_id"`
	ShopName        string    `json:"shop_name"`
	Quantity        int       `json:"quantity"`
	Subtotal        float64   `json:"subtotal"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CalculateSubtotal calculates subtotal for cart item
func (item *CartItemWithProduct) CalculateSubtotal() {
	item.Subtotal = item.ProductPrice * float64(item.Quantity)
}

// CalculateTotals calculates total items and subtotal for cart
func (cart *CartResponse) CalculateTotals() {
	cart.ItemCount = len(cart.Items)
	cart.Subtotal = 0
	for _, item := range cart.Items {
		cart.Subtotal += item.Subtotal
	}
}
