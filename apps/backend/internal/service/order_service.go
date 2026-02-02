package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	cartRepo    *repository.CartRepository
	productRepo *repository.ProductRepository
}

func NewOrderService(
	orderRepo *repository.OrderRepository,
	cartRepo *repository.CartRepository,
	productRepo *repository.ProductRepository,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

// CreateOrderFromCart creates an order from user's cart
func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID uuid.UUID, req *model.CreateOrderRequest) (*model.OrderResponse, error) {
	// Get user's cart
	cart, err := s.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cart not found")
	}

	// Get cart items with product details
	cartItems, err := s.cartRepo.GetCartWithItems(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Validate stock availability for all items
	for _, item := range cartItems {
		product, err := s.productRepo.GetByID(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found", item.ProductName)
		}

		if !product.IsActive {
			return nil, fmt.Errorf("product %s is no longer available", item.ProductName)
		}

		if product.StockQuantity < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s: only %d available", item.ProductName, product.StockQuantity)
		}
	}

	// Create shipping address
	shippingAddress := &model.Address{
		ID:           uuid.New(),
		UserID:       userID,
		FullName:     req.ShippingAddress.FullName,
		Phone:        req.ShippingAddress.Phone,
		AddressLine1: req.ShippingAddress.AddressLine1,
		AddressLine2: req.ShippingAddress.AddressLine2,
		City:         req.ShippingAddress.City,
		State:        req.ShippingAddress.State,
		PostalCode:   req.ShippingAddress.PostalCode,
		Country:      req.ShippingAddress.Country,
		IsDefault:    false,
		AddressType:  "shipping",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.orderRepo.CreateAddress(ctx, shippingAddress); err != nil {
		return nil, fmt.Errorf("failed to create shipping address: %w", err)
	}

	// Create billing address
	var billingAddressID *uuid.UUID
	if req.UseSameAddress {
		billingAddressID = &shippingAddress.ID
	} else if req.BillingAddress != nil {
		billingAddress := &model.Address{
			ID:           uuid.New(),
			UserID:       userID,
			FullName:     req.BillingAddress.FullName,
			Phone:        req.BillingAddress.Phone,
			AddressLine1: req.BillingAddress.AddressLine1,
			AddressLine2: req.BillingAddress.AddressLine2,
			City:         req.BillingAddress.City,
			State:        req.BillingAddress.State,
			PostalCode:   req.BillingAddress.PostalCode,
			Country:      req.BillingAddress.Country,
			IsDefault:    false,
			AddressType:  "billing",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.orderRepo.CreateAddress(ctx, billingAddress); err != nil {
			return nil, fmt.Errorf("failed to create billing address: %w", err)
		}
		billingAddressID = &billingAddress.ID
	}

	// Calculate totals
	var subtotal float64
	for _, item := range cartItems {
		subtotal += item.Subtotal
	}

	shippingCost := 0.0 // TODO: Calculate based on location/weight
	tax := 0.0          // TODO: Calculate based on region
	discount := 0.0
	total := subtotal + shippingCost + tax - discount

	// Create order
	order := &model.Order{
		ID:                uuid.New(),
		UserID:            userID,
		OrderNumber:       model.GenerateOrderNumber(),
		Status:            model.OrderStatusPending,
		ShippingAddressID: &shippingAddress.ID,
		BillingAddressID:  billingAddressID,
		Subtotal:          subtotal,
		ShippingCost:      shippingCost,
		Tax:               tax,
		Discount:          discount,
		Total:             total,
		PaymentMethod:     &req.PaymentMethod,
		PaymentStatus:     model.PaymentStatusPending,
		Notes:             req.Notes,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Create order items and reduce stock
	var orderItems []model.OrderItem
	for _, cartItem := range cartItems {
		orderItem := model.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ShopID:      cartItem.ShopID,
			ProductName: cartItem.ProductName,
			Quantity:    cartItem.Quantity,
			UnitPrice:   cartItem.ProductPrice,
			Subtotal:    cartItem.Subtotal,
			CreatedAt:   time.Now(),
		}
		orderItems = append(orderItems, orderItem)

		// Reduce product stock
		if err := s.productRepo.ReduceStock(ctx, cartItem.ProductID, cartItem.Quantity); err != nil {
			return nil, fmt.Errorf("failed to reduce stock for %s: %w", cartItem.ProductName, err)
		}
	}

	if err := s.orderRepo.CreateOrderItems(ctx, orderItems); err != nil {
		return nil, fmt.Errorf("failed to create order items: %w", err)
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(ctx, cart.ID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Get full order details for response
	return s.GetOrderByID(ctx, order.ID, userID)
}

// GetOrderByID retrieves order by ID with authorization check
func (s *OrderService) GetOrderByID(ctx context.Context, orderID, userID uuid.UUID) (*model.OrderResponse, error) {
	// Verify ownership
	owned, err := s.orderRepo.VerifyOrderOwnership(ctx, orderID, userID)
	if err != nil || !owned {
		return nil, fmt.Errorf("order not found or unauthorized")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	// Get order items
	items, err := s.orderRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	// Get addresses
	var shippingAddress *model.Address
	var billingAddress *model.Address

	if order.ShippingAddressID != nil {
		shippingAddress, err = s.orderRepo.GetAddressByID(ctx, *order.ShippingAddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get shipping address: %w", err)
		}
	}

	if order.BillingAddressID != nil {
		billingAddress, err = s.orderRepo.GetAddressByID(ctx, *order.BillingAddressID)
		if err != nil {
			return nil, fmt.Errorf("failed to get billing address: %w", err)
		}
	}

	return &model.OrderResponse{
		ID:              order.ID,
		UserID:          order.UserID,
		OrderNumber:     order.OrderNumber,
		Status:          order.Status,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		Items:           items,
		Subtotal:        order.Subtotal,
		ShippingCost:    order.ShippingCost,
		Tax:             order.Tax,
		Discount:        order.Discount,
		Total:           order.Total,
		PaymentMethod:   order.PaymentMethod,
		PaymentStatus:   order.PaymentStatus,
		Notes:           order.Notes,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
		ConfirmedAt:     order.ConfirmedAt,
		ShippedAt:       order.ShippedAt,
		DeliveredAt:     order.DeliveredAt,
	}, nil
}

// GetUserOrders retrieves all orders for a user
func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	return s.orderRepo.GetByUserID(ctx, userID)
}

// GetVendorOrders retrieves orders for vendor's shop
func (s *OrderService) GetVendorOrders(ctx context.Context, shopID uuid.UUID) ([]*model.Order, error) {
	return s.orderRepo.GetByShopID(ctx, shopID)
}

// UpdateOrderStatus updates order status (vendor/admin only)
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	return s.orderRepo.UpdateStatusWithTimestamp(ctx, orderID, status)
}

// CancelOrder cancels an order and restores stock
func (s *OrderService) CancelOrder(ctx context.Context, orderID, userID uuid.UUID) error {
	// Verify ownership
	owned, err := s.orderRepo.VerifyOrderOwnership(ctx, orderID, userID)
	if err != nil || !owned {
		return fmt.Errorf("order not found or unauthorized")
	}

	// Get order
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	// Check if order can be cancelled
	if order.Status != model.OrderStatusPending && order.Status != model.OrderStatusConfirmed {
		return fmt.Errorf("order cannot be cancelled in current status: %s", order.Status)
	}

	// Get order items
	items, err := s.orderRepo.GetOrderItems(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get order items: %w", err)
	}

	// Restore stock for each item
	for _, item := range items {
		if err := s.productRepo.IncreaseStock(ctx, item.ProductID, item.Quantity); err != nil {
			return fmt.Errorf("failed to restore stock for %s: %w", item.ProductName, err)
		}
	}

	// Update order status to cancelled
	return s.orderRepo.UpdateStatus(ctx, orderID, model.OrderStatusCancelled)
}
