package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
)

type OrderRepository struct {
	db *database.Database
}

func NewOrderRepository(db *database.Database) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(ctx context.Context, order *model.Order) error {
	query := `
		INSERT INTO orders (
			id, user_id, order_number, status, shipping_address_id, billing_address_id,
			subtotal, shipping_cost, tax, discount, total, payment_method, payment_status,
			notes, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.ShippingAddressID,
		order.BillingAddressID,
		order.Subtotal,
		order.ShippingCost,
		order.Tax,
		order.Discount,
		order.Total,
		order.PaymentMethod,
		order.PaymentStatus,
		order.Notes,
		order.CreatedAt,
		order.UpdatedAt,
	)

	return err
}

// CreateOrderItems creates order items in batch
func (r *OrderRepository) CreateOrderItems(ctx context.Context, items []model.OrderItem) error {
	query := `
		INSERT INTO order_items (
			id, order_id, product_id, shop_id, product_name, product_sku,
			quantity, unit_price, subtotal, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for _, item := range items {
		_, err := r.db.Pool.Exec(ctx, query,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.ShopID,
			item.ProductName,
			item.ProductSKU,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
			item.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves order by ID with items
func (r *OrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	var order model.Order
	query := `
		SELECT id, user_id, order_number, status, shipping_address_id, billing_address_id,
		       subtotal, shipping_cost, tax, discount, total, payment_method, payment_status,
		       notes, created_at, updated_at, confirmed_at, shipped_at, delivered_at
		FROM orders
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.ShippingAddressID,
		&order.BillingAddressID,
		&order.Subtotal,
		&order.ShippingCost,
		&order.Tax,
		&order.Discount,
		&order.Total,
		&order.PaymentMethod,
		&order.PaymentStatus,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.ConfirmedAt,
		&order.ShippedAt,
		&order.DeliveredAt,
	)

	return &order, err
}

// GetOrderItems retrieves all items for an order
func (r *OrderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]model.OrderItemWithDetails, error) {
	query := `
		SELECT 
			oi.id,
			oi.order_id,
			oi.product_id,
			oi.product_name,
			p.image_url as product_image_url,
			oi.shop_id,
			s.shop_name,
			oi.quantity,
			oi.unit_price,
			oi.subtotal,
			oi.created_at
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		LEFT JOIN shops s ON oi.shop_id = s.id
		WHERE oi.order_id = $1
		ORDER BY oi.created_at ASC
	`

	rows, err := r.db.Pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.OrderItemWithDetails
	for rows.Next() {
		var item model.OrderItemWithDetails
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductImageURL,
			&item.ShopID,
			&item.ShopName,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetByUserID retrieves all orders for a user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Order, error) {
	query := `
		SELECT id, user_id, order_number, status, shipping_address_id, billing_address_id,
		       subtotal, shipping_cost, tax, discount, total, payment_method, payment_status,
		       notes, created_at, updated_at, confirmed_at, shipped_at, delivered_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.ShippingAddressID,
			&order.BillingAddressID,
			&order.Subtotal,
			&order.ShippingCost,
			&order.Tax,
			&order.Discount,
			&order.Total,
			&order.PaymentMethod,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.ConfirmedAt,
			&order.ShippedAt,
			&order.DeliveredAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// GetByShopID retrieves orders containing products from a specific shop
func (r *OrderRepository) GetByShopID(ctx context.Context, shopID uuid.UUID) ([]*model.Order, error) {
	query := `
		SELECT DISTINCT o.id, o.user_id, o.order_number, o.status, o.shipping_address_id, o.billing_address_id,
		       o.subtotal, o.shipping_cost, o.tax, o.discount, o.total, o.payment_method, o.payment_status,
		       o.notes, o.created_at, o.updated_at, o.confirmed_at, o.shipped_at, o.delivered_at
		FROM orders o
		INNER JOIN order_items oi ON o.id = oi.order_id
		WHERE oi.shop_id = $1
		ORDER BY o.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.Status,
			&order.ShippingAddressID,
			&order.BillingAddressID,
			&order.Subtotal,
			&order.ShippingCost,
			&order.Tax,
			&order.Discount,
			&order.Total,
			&order.PaymentMethod,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.ConfirmedAt,
			&order.ShippedAt,
			&order.DeliveredAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, rows.Err()
}

// UpdateStatus updates order status
func (r *OrderRepository) UpdateStatus(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// UpdateStatusWithTimestamp updates order status and sets appropriate timestamp
func (r *OrderRepository) UpdateStatusWithTimestamp(ctx context.Context, orderID uuid.UUID, status model.OrderStatus) error {
	var query string

	switch status {
	case model.OrderStatusConfirmed:
		query = `UPDATE orders SET status = $1, confirmed_at = NOW(), updated_at = NOW() WHERE id = $2`
	case model.OrderStatusShipped:
		query = `UPDATE orders SET status = $1, shipped_at = NOW(), updated_at = NOW() WHERE id = $2`
	case model.OrderStatusDelivered:
		query = `UPDATE orders SET status = $1, delivered_at = NOW(), updated_at = NOW() WHERE id = $2`
	default:
		query = `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	}

	result, err := r.db.Pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// UpdatePaymentStatus updates the payment status of an order
func (r *OrderRepository) UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, status model.PaymentStatus) error {
	query := `UPDATE orders SET payment_status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}

// CreateAddress creates a new address
func (r *OrderRepository) CreateAddress(ctx context.Context, address *model.Address) error {
	query := `
		INSERT INTO addresses (
			id, user_id, full_name, phone, address_line1, address_line2,
			city, state, postal_code, country, is_default, address_type,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		address.ID,
		address.UserID,
		address.FullName,
		address.Phone,
		address.AddressLine1,
		address.AddressLine2,
		address.City,
		address.State,
		address.PostalCode,
		address.Country,
		address.IsDefault,
		address.AddressType,
		address.CreatedAt,
		address.UpdatedAt,
	)

	return err
}

// GetAddressByID retrieves address by ID
func (r *OrderRepository) GetAddressByID(ctx context.Context, addressID uuid.UUID) (*model.Address, error) {
	var address model.Address
	query := `
		SELECT id, user_id, full_name, phone, address_line1, address_line2,
		       city, state, postal_code, country, is_default, address_type,
		       created_at, updated_at
		FROM addresses
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, addressID).Scan(
		&address.ID,
		&address.UserID,
		&address.FullName,
		&address.Phone,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.City,
		&address.State,
		&address.PostalCode,
		&address.Country,
		&address.IsDefault,
		&address.AddressType,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	return &address, err
}

// VerifyOrderOwnership checks if order belongs to user
func (r *OrderRepository) VerifyOrderOwnership(ctx context.Context, orderID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1 AND user_id = $2)`
	err := r.db.Pool.QueryRow(ctx, query, orderID, userID).Scan(&exists)
	return exists, err
}
