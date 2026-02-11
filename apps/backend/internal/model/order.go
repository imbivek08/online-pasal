package model

import (
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type PaymentStatus string

const (
	OrderStatusPending    OrderStatus = "pending"
	OrderStatusConfirmed  OrderStatus = "confirmed"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
	OrderStatusRefunded   OrderStatus = "refunded"
)

const (
	PaymentStatusPending  PaymentStatus = "pending"
	PaymentStatusPaid     PaymentStatus = "paid"
	PaymentStatusFailed   PaymentStatus = "failed"
	PaymentStatusRefunded PaymentStatus = "refunded"
)

// Order represents a customer order
type Order struct {
	ID                uuid.UUID     `json:"id" db:"id"`
	UserID            uuid.UUID     `json:"user_id" db:"user_id"`
	OrderNumber       string        `json:"order_number" db:"order_number"`
	Status            OrderStatus   `json:"status" db:"status"`
	ShippingAddressID *uuid.UUID    `json:"shipping_address_id,omitempty" db:"shipping_address_id"`
	BillingAddressID  *uuid.UUID    `json:"billing_address_id,omitempty" db:"billing_address_id"`
	Subtotal          float64       `json:"subtotal" db:"subtotal"`
	ShippingCost      float64       `json:"shipping_cost" db:"shipping_cost"`
	Tax               float64       `json:"tax" db:"tax"`
	Discount          float64       `json:"discount" db:"discount"`
	Total             float64       `json:"total" db:"total"`
	PaymentMethod     *string       `json:"payment_method,omitempty" db:"payment_method"`
	PaymentStatus     PaymentStatus `json:"payment_status" db:"payment_status"`
	StripeSessionID   *string       `json:"stripe_session_id,omitempty" db:"stripe_session_id"`
	Notes             *string       `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at" db:"updated_at"`
	ConfirmedAt       *time.Time    `json:"confirmed_at,omitempty" db:"confirmed_at"`
	ShippedAt         *time.Time    `json:"shipped_at,omitempty" db:"shipped_at"`
	DeliveredAt       *time.Time    `json:"delivered_at,omitempty" db:"delivered_at"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID          uuid.UUID `json:"id" db:"id"`
	OrderID     uuid.UUID `json:"order_id" db:"order_id"`
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	ShopID      uuid.UUID `json:"shop_id" db:"shop_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	ProductSKU  *string   `json:"product_sku,omitempty" db:"product_sku"`
	Quantity    int       `json:"quantity" db:"quantity"`
	UnitPrice   float64   `json:"unit_price" db:"unit_price"`
	Subtotal    float64   `json:"subtotal" db:"subtotal"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// OrderItemWithDetails includes product and shop information
type OrderItemWithDetails struct {
	ID              uuid.UUID `json:"id" db:"id"`
	OrderID         uuid.UUID `json:"order_id" db:"order_id"`
	ProductID       uuid.UUID `json:"product_id" db:"product_id"`
	ProductName     string    `json:"product_name" db:"product_name"`
	ProductImageURL *string   `json:"product_image_url,omitempty" db:"product_image_url"`
	ShopID          uuid.UUID `json:"shop_id" db:"shop_id"`
	ShopName        string    `json:"shop_name" db:"shop_name"`
	Quantity        int       `json:"quantity" db:"quantity"`
	UnitPrice       float64   `json:"unit_price" db:"unit_price"`
	Subtotal        float64   `json:"subtotal" db:"subtotal"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// Address represents a shipping or billing address
type Address struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	FullName     string    `json:"full_name" db:"full_name"`
	Phone        string    `json:"phone" db:"phone"`
	AddressLine1 string    `json:"address_line1" db:"address_line1"`
	AddressLine2 *string   `json:"address_line2,omitempty" db:"address_line2"`
	City         string    `json:"city" db:"city"`
	State        *string   `json:"state,omitempty" db:"state"`
	PostalCode   *string   `json:"postal_code,omitempty" db:"postal_code"`
	Country      string    `json:"country" db:"country"`
	IsDefault    bool      `json:"is_default" db:"is_default"`
	AddressType  string    `json:"address_type" db:"address_type"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateOrderRequest represents request to create an order from cart
type CreateOrderRequest struct {
	ShippingAddressID *uuid.UUID    `json:"shipping_address_id,omitempty"` // Use existing saved address
	ShippingAddress   *AddressInput `json:"shipping_address,omitempty"`    // Or provide new address
	BillingAddress    *AddressInput `json:"billing_address,omitempty"`
	PaymentMethod     string        `json:"payment_method" validate:"required"`
	UseSameAddress    bool          `json:"use_same_address"` // Use shipping as billing
	Notes             *string       `json:"notes,omitempty"`
}

// AddressInput represents address input for order creation and address management
type AddressInput struct {
	FullName     string  `json:"full_name" validate:"required"`
	Phone        string  `json:"phone" validate:"required"`
	AddressLine1 string  `json:"address_line1" validate:"required"`
	AddressLine2 *string `json:"address_line2,omitempty"`
	City         string  `json:"city" validate:"required"`
	State        *string `json:"state,omitempty"`
	PostalCode   *string `json:"postal_code,omitempty"`
	Country      string  `json:"country" validate:"required"`
	IsDefault    bool    `json:"is_default"`
}

// UpdateOrderStatusRequest represents request to update order status
type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}

// OrderResponse represents order with items
type OrderResponse struct {
	ID              uuid.UUID              `json:"id"`
	UserID          uuid.UUID              `json:"user_id"`
	OrderNumber     string                 `json:"order_number"`
	Status          OrderStatus            `json:"status"`
	ShippingAddress *Address               `json:"shipping_address,omitempty"`
	BillingAddress  *Address               `json:"billing_address,omitempty"`
	Items           []OrderItemWithDetails `json:"items"`
	Subtotal        float64                `json:"subtotal"`
	ShippingCost    float64                `json:"shipping_cost"`
	Tax             float64                `json:"tax"`
	Discount        float64                `json:"discount"`
	Total           float64                `json:"total"`
	PaymentMethod   *string                `json:"payment_method,omitempty"`
	PaymentStatus   PaymentStatus          `json:"payment_status"`
	Notes           *string                `json:"notes,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	ConfirmedAt     *time.Time             `json:"confirmed_at,omitempty"`
	ShippedAt       *time.Time             `json:"shipped_at,omitempty"`
	DeliveredAt     *time.Time             `json:"delivered_at,omitempty"`
}

// OrderSummary represents a simplified order for lists
type OrderSummary struct {
	ID          uuid.UUID   `json:"id"`
	OrderNumber string      `json:"order_number"`
	Status      OrderStatus `json:"status"`
	ItemCount   int         `json:"item_count"`
	Total       float64     `json:"total"`
	CreatedAt   time.Time   `json:"created_at"`
}

// StripeSessionStatus represents status info from a Stripe checkout session
type StripeSessionStatus struct {
	SessionID     string    `json:"session_id"`
	PaymentStatus string    `json:"payment_status"`
	OrderID       uuid.UUID `json:"order_id"`
	OrderNumber   string    `json:"order_number"`
}

// CreateCheckoutResponse is returned when a Stripe checkout order is created
type CreateCheckoutResponse struct {
	Order       *OrderResponse `json:"order"`
	CheckoutURL string         `json:"checkout_url,omitempty"`
}

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	// Format: ORD-YYYYMMDD-XXXXX
	now := time.Now()
	timestamp := now.Format("20060102")
	random := uuid.New().String()[:8]
	return "ORD-" + timestamp + "-" + random
}
