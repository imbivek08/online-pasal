package repository

import (
	"context"

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
		INSERT INTO products (id, vendor_id, name, description, price, stock_quantity, category, image_url, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Pool.Exec(ctx, query,
		product.ID,
		product.VendorID,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
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
		SELECT id, vendor_id, name, description, price, stock_quantity, category, image_url, is_active, created_at, updated_at
		FROM products
		WHERE id = $1
	`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.VendorID,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.StockQuantity,
		&product.Category,
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

// GetAll retrieves all products with optional filters
func (r *ProductRepository) GetAll(ctx context.Context, filters map[string]interface{}) ([]*model.Product, error) {
	var products []*model.Product
	query := `
		SELECT id, vendor_id, name, description, price, stock_quantity, category, image_url, is_active, created_at, updated_at
		FROM products
		WHERE is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.VendorID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
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

// GetByVendorID retrieves all products for a specific vendor
func (r *ProductRepository) GetByVendorID(ctx context.Context, vendorID uuid.UUID) ([]*model.Product, error) {
	var products []*model.Product
	query := `
		SELECT id, vendor_id, name, description, price, stock_quantity, category, image_url, is_active, created_at, updated_at
		FROM products
		WHERE vendor_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Pool.Query(ctx, query, vendorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ID,
			&product.VendorID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.StockQuantity,
			&product.Category,
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
		SET name = $1, description = $2, price = $3, stock_quantity = $4, category = $5, image_url = $6, is_active = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Pool.Exec(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.StockQuantity,
		product.Category,
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
