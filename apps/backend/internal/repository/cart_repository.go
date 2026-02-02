package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
)

type CartRepository struct {
	db *database.Database
}

func NewCartRepository(db *database.Database) *CartRepository {
	return &CartRepository{db: db}
}

// GetOrCreateCart gets existing cart or creates new one for user
func (r *CartRepository) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	// Try to get existing cart
	cart, err := r.GetCartByUserID(ctx, userID)
	if err == nil {
		return cart, nil
	}

	// Create new cart if not found
	cart = &model.Cart{
		ID:     uuid.New(),
		UserID: &userID,
	}

	query := `
		INSERT INTO carts (id, user_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, user_id, session_id, created_at, updated_at, expires_at
	`

	err = r.db.Pool.QueryRow(ctx, query, cart.ID, cart.UserID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.SessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.ExpiresAt,
	)

	return cart, err
}

// GetCartByUserID retrieves cart by user ID
func (r *CartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*model.Cart, error) {
	var cart model.Cart
	query := `
		SELECT id, user_id, session_id, created_at, updated_at, expires_at
		FROM carts
		WHERE user_id = $1
		LIMIT 1
	`

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&cart.ID,
		&cart.UserID,
		&cart.SessionID,
		&cart.CreatedAt,
		&cart.UpdatedAt,
		&cart.ExpiresAt,
	)

	return &cart, err
}

// AddItem adds or updates item in cart
func (r *CartRepository) AddItem(ctx context.Context, cartID, productID uuid.UUID, quantity int) (*model.CartItem, error) {
	var item model.CartItem

	query := `
		INSERT INTO cart_items (id, cart_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (cart_id, product_id) 
		DO UPDATE SET 
			quantity = cart_items.quantity + EXCLUDED.quantity,
			updated_at = NOW()
		RETURNING id, cart_id, product_id, quantity, created_at, updated_at
	`

	err := r.db.Pool.QueryRow(ctx, query, uuid.New(), cartID, productID, quantity).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	return &item, err
}

// UpdateItemQuantity updates cart item quantity
func (r *CartRepository) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := r.db.Pool.Exec(ctx, query, quantity, itemID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// RemoveItem removes item from cart
func (r *CartRepository) RemoveItem(ctx context.Context, itemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`

	result, err := r.db.Pool.Exec(ctx, query, itemID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}

// GetCartWithItems retrieves cart with all items and product details
func (r *CartRepository) GetCartWithItems(ctx context.Context, cartID uuid.UUID) ([]model.CartItemWithProduct, error) {
	query := `
		SELECT 
			ci.id,
			ci.cart_id,
			ci.product_id,
			p.name as product_name,
			p.price as product_price,
			p.image_url as product_image_url,
			p.stock_quantity,
			p.is_active,
			p.shop_id,
			s.shop_name,
			ci.quantity,
			ci.created_at,
			ci.updated_at
		FROM cart_items ci
		INNER JOIN products p ON ci.product_id = p.id
		INNER JOIN shops s ON p.shop_id = s.id
		WHERE ci.cart_id = $1
		ORDER BY ci.created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.CartItemWithProduct
	for rows.Next() {
		var item model.CartItemWithProduct
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.ProductName,
			&item.ProductPrice,
			&item.ProductImageURL,
			&item.StockQuantity,
			&item.IsActive,
			&item.ShopID,
			&item.ShopName,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Calculate subtotal
		item.CalculateSubtotal()
		items = append(items, item)
	}

	return items, rows.Err()
}

// ClearCart removes all items from cart
func (r *CartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.Pool.Exec(ctx, query, cartID)
	return err
}

// GetCartItemCount returns the total number of items in cart
func (r *CartRepository) GetCartItemCount(ctx context.Context, cartID uuid.UUID) (int, error) {
	var count int
	query := `SELECT COALESCE(SUM(quantity), 0) FROM cart_items WHERE cart_id = $1`
	err := r.db.Pool.QueryRow(ctx, query, cartID).Scan(&count)
	return count, err
}

// GetCartItemByID retrieves a specific cart item
func (r *CartRepository) GetCartItemByID(ctx context.Context, itemID uuid.UUID) (*model.CartItem, error) {
	var item model.CartItem
	query := `
		SELECT id, cart_id, product_id, quantity, created_at, updated_at
		FROM cart_items
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, itemID).Scan(
		&item.ID,
		&item.CartID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	return &item, err
}

// VerifyCartOwnership checks if cart belongs to user
func (r *CartRepository) VerifyCartOwnership(ctx context.Context, cartID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM carts WHERE id = $1 AND user_id = $2)`
	err := r.db.Pool.QueryRow(ctx, query, cartID, userID).Scan(&exists)
	return exists, err
}
