package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
)

type ProductRepository struct {
	db *database.Database
}

func NewProductRepository(db *database.Database) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (id, shop_id, category_id, name, description, price, stock_quantity, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		product.ID,
		product.ShopID,
		product.CategoryID,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.ImageURL,
		product.IsActive,
		product.CreatedAt,
		product.UpdatedAt,
	)
	return err
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error) {
	var product model.Product
	query := `
		SELECT id, shop_id, category_id, name, description, price, stock_quantity, image_url, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.ShopID,
		&product.CategoryID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// GetShopByProductID retrieves shop information for a product
func (r *ProductRepository) GetShopByProductID(ctx context.Context, productID uuid.UUID) (*model.Shop, error) {
	var shop model.Shop
	query := `
		SELECT s.id, s.vendor_id, s.shop_name as name, s.slug, s.description, s.logo_url, s.banner_url, 
		       s.contact_email, s.contact_phone, s.address, s.city, s.state, s.country, s.postal_code,
		       s.is_active, s.is_verified, s.created_at, s.updated_at
		FROM shops s
		INNER JOIN products p ON s.id = p.shop_id
		WHERE p.id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, productID).Scan(
		&shop.ID,
		&shop.VendorID,
		&shop.Name,
		&shop.Slug,
		&shop.Description,
		&shop.LogoURL,
		&shop.BannerURL,
		&shop.Email,
		&shop.Phone,
		&shop.Address,
		&shop.City,
		&shop.State,
		&shop.Country,
		&shop.PostalCode,
		&shop.IsActive,
		&shop.IsVerified,
		&shop.CreatedAt,
		&shop.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &shop, nil
}

// GetAll retrieves all products with optional filters
func (r *ProductRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*model.Product, error) {
	var products []*model.Product
	query := `
		SELECT id, shop_id, category_id, name, description, price, stock_quantity, image_url, is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
	`
	args := []interface{}{}
	argIndex := 1

	// Search filter (matches name or description)
	if search, ok := filters["search"].(string); ok && search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Min price filter
	if minPrice, ok := filters["min_price"].(float64); ok {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, minPrice)
		argIndex++
	}

	// Max price filter
	if maxPrice, ok := filters["max_price"].(float64); ok {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, maxPrice)
		argIndex++
	}

	// Sort
	sortBy, _ := filters["sort_by"].(string)
	switch sortBy {
	case "price_asc":
		query += " ORDER BY price ASC"
	case "price_desc":
		query += " ORDER BY price DESC"
	case "name_asc":
		query += " ORDER BY name ASC"
	default:
		query += " ORDER BY created_at DESC"
	}

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

// GetByShopID retrieves all products for a specific shop
func (r *ProductRepository) GetByShopID(ctx context.Context, shopID uuid.UUID) ([]*model.Product, error) {
	var products []*model.Product
	query := `
		SELECT id, shop_id, name, description, price, stock_quantity, category_id, image_url, is_active, created_at, updated_at
		FROM products
		WHERE shop_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.ShopID,
			&product.CategoryID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.ImageURL,
			&product.IsActive,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, rows.Err()
}

// Update updates a product
func (r *ProductRepository) Update(ctx context.Context, product *model.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock_quantity = $4, category_id = $5, image_url = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Pool.Exec(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.CategoryID,
		product.ImageURL,
		product.IsActive,
		product.UpdatedAt,
		product.ID,
	)
	return err
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

// ReduceStock reduces product stock quantity
func (r *ProductRepository) ReduceStock(ctx context.Context, productID uuid.UUID, quantity int) error {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity - $1, updated_at = NOW()
		WHERE id = $2 AND stock_quantity >= $1
	`
	result, err := r.db.Pool.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock or product not found")
	}

	return nil
}

// IncreaseStock increases product stock quantity
func (r *ProductRepository) IncreaseStock(ctx context.Context, productID uuid.UUID, quantity int) error {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity + $1, updated_at = NOW()
		WHERE id = $2
	`
	result, err := r.db.Pool.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}
