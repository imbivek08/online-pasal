package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/imbivek08/hamropasal/internal/model"
)

type ShopRepository struct {
	db *pgxpool.Pool
}

func NewShopRepository(db *pgxpool.Pool) *ShopRepository {
	return &ShopRepository{db: db}
}

// Create creates a new shop
func (r *ShopRepository) Create(ctx context.Context, shop *model.Shop) error {
	query := `
		INSERT INTO shops (id, vendor_id, shop_name, slug, description, logo_url, banner_url, 
		                   address, city, state, country, postal_code, contact_phone, contact_email, 
		                   is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		shop.ID,
		shop.VendorID,
		shop.Name,
		shop.Slug,
		shop.Description,
		shop.LogoURL,
		shop.BannerURL,
		shop.Address,
		shop.City,
		shop.State,
		shop.Country,
		shop.PostalCode,
		shop.Phone,
		shop.Email,
		shop.IsActive,
		shop.IsVerified,
	).Scan(&shop.CreatedAt, &shop.UpdatedAt)

	return err
}

// GetByID retrieves a shop by ID
func (r *ShopRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Shop, error) {
	query := `
		SELECT id, vendor_id, shop_name, slug, description, logo_url, banner_url,
		       address, city, state, country, postal_code, contact_phone, contact_email,
		       is_active, is_verified, created_at, updated_at
		FROM shops
		WHERE id = $1
	`

	shop := &model.Shop{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&shop.ID,
		&shop.VendorID,
		&shop.Name,
		&shop.Slug,
		&shop.Description,
		&shop.LogoURL,
		&shop.BannerURL,
		&shop.Address,
		&shop.City,
		&shop.State,
		&shop.Country,
		&shop.PostalCode,
		&shop.Phone,
		&shop.Email,
		&shop.IsActive,
		&shop.IsVerified,
		&shop.CreatedAt,
		&shop.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return shop, nil
}

// GetBySlug retrieves a shop by slug
func (r *ShopRepository) GetBySlug(ctx context.Context, slug string) (*model.Shop, error) {
	query := `
		SELECT id, vendor_id, shop_name, slug, description, logo_url, banner_url,
		       address, city, state, country, postal_code, contact_phone, contact_email,
		       is_active, is_verified, created_at, updated_at
		FROM shops
		WHERE slug = $1
	`

	shop := &model.Shop{}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&shop.ID,
		&shop.VendorID,
		&shop.Name,
		&shop.Slug,
		&shop.Description,
		&shop.LogoURL,
		&shop.BannerURL,
		&shop.Address,
		&shop.City,
		&shop.State,
		&shop.Country,
		&shop.PostalCode,
		&shop.Phone,
		&shop.Email,
		&shop.IsActive,
		&shop.IsVerified,
		&shop.CreatedAt,
		&shop.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return shop, nil
}

// GetByVendorID retrieves a shop by vendor ID
func (r *ShopRepository) GetByVendorID(ctx context.Context, vendorID uuid.UUID) (*model.Shop, error) {
	query := `
		SELECT id, vendor_id, shop_name, slug, description, logo_url, banner_url,
		       address, city, state, country, postal_code, contact_phone, contact_email,
		       is_active, is_verified, created_at, updated_at
		FROM shops
		WHERE vendor_id = $1
	`

	shop := &model.Shop{}
	err := r.db.QueryRow(ctx, query, vendorID).Scan(
		&shop.ID,
		&shop.VendorID,
		&shop.Name,
		&shop.Slug,
		&shop.Description,
		&shop.LogoURL,
		&shop.BannerURL,
		&shop.Address,
		&shop.City,
		&shop.State,
		&shop.Country,
		&shop.PostalCode,
		&shop.Phone,
		&shop.Email,
		&shop.IsActive,
		&shop.IsVerified,
		&shop.CreatedAt,
		&shop.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return shop, nil
}

// List retrieves shops with pagination and filters
func (r *ShopRepository) List(ctx context.Context, page, pageSize int, search string, activeOnly bool) ([]model.Shop, int, error) {
	offset := (page - 1) * pageSize

	// Build query with filters
	whereConditions := []string{}
	args := []interface{}{}
	argCounter := 1

	if activeOnly {
		whereConditions = append(whereConditions, fmt.Sprintf("is_active = $%d", argCounter))
		args = append(args, true)
		argCounter++
	}

	if search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argCounter, argCounter))
		args = append(args, "%"+search+"%")
		argCounter++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM shops %s", whereClause)
	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get shops
	args = append(args, pageSize, offset)
	query := fmt.Sprintf(`
		SELECT id, vendor_id, shop_name, slug, description, logo_url, banner_url,
		       address, city, state, country, postal_code, contact_phone, contact_email,
		       is_active, is_verified, created_at, updated_at
		FROM shops
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCounter, argCounter+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	shops := []model.Shop{}
	for rows.Next() {
		var shop model.Shop
		err := rows.Scan(
			&shop.ID,
			&shop.VendorID,
			&shop.Name,
			&shop.Slug,
			&shop.Description,
			&shop.LogoURL,
			&shop.BannerURL,
			&shop.Address,
			&shop.City,
			&shop.State,
			&shop.Country,
			&shop.PostalCode,
			&shop.Phone,
			&shop.Email,
			&shop.IsActive,
			&shop.IsVerified,
			&shop.CreatedAt,
			&shop.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		shops = append(shops, shop)
	}

	return shops, total, nil
}

// Update updates a shop
func (r *ShopRepository) Update(ctx context.Context, shop *model.Shop) error {
	query := `
		UPDATE shops
		SET name = $1, description = $2, logo_url = $3, banner_url = $4,
		    address = $5, city = $6, state = $7, country = $8, postal_code = $9,
		    phone = $10, email = $11, updated_at = NOW()
		WHERE id = $12
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		shop.Name,
		shop.Description,
		shop.LogoURL,
		shop.BannerURL,
		shop.Address,
		shop.City,
		shop.State,
		shop.Country,
		shop.PostalCode,
		shop.Phone,
		shop.Email,
		shop.ID,
	).Scan(&shop.UpdatedAt)

	return err
}

// UpdateStatus updates shop's active status
func (r *ShopRepository) UpdateStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	query := `UPDATE shops SET is_active = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, isActive, id)
	return err
}

// UpdateVerification updates shop's verification status
func (r *ShopRepository) UpdateVerification(ctx context.Context, id uuid.UUID, isVerified bool) error {
	query := `UPDATE shops SET is_verified = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, isVerified, id)
	return err
}

// Delete soft deletes a shop
func (r *ShopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE shops SET is_active = false, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetStats retrieves shop statistics
func (r *ShopRepository) GetStats(ctx context.Context, shopID uuid.UUID) (*model.ShopWithStats, error) {
	query := `
		SELECT 
			s.id, s.vendor_id, s.name, s.slug, s.description, s.logo_url, s.banner_url,
			s.address, s.city, s.state, s.country, s.postal_code, s.phone, s.email,
			s.is_active, s.is_verified, s.created_at, s.updated_at,
			COUNT(DISTINCT p.id) as total_products,
			COUNT(DISTINCT oi.order_id) as total_orders,
			COALESCE(SUM(oi.quantity * oi.unit_price), 0) as total_revenue,
			COALESCE(AVG(r.rating), 0) as average_rating
		FROM shops s
		LEFT JOIN products p ON p.shop_id = s.id
		LEFT JOIN order_items oi ON oi.product_id = p.id
		LEFT JOIN reviews r ON r.product_id = p.id
		WHERE s.id = $1
		GROUP BY s.id
	`

	stats := &model.ShopWithStats{}
	err := r.db.QueryRow(ctx, query, shopID).Scan(
		&stats.ID,
		&stats.VendorID,
		&stats.Name,
		&stats.Slug,
		&stats.Description,
		&stats.LogoURL,
		&stats.BannerURL,
		&stats.Address,
		&stats.City,
		&stats.State,
		&stats.Country,
		&stats.PostalCode,
		&stats.Phone,
		&stats.Email,
		&stats.IsActive,
		&stats.IsVerified,
		&stats.CreatedAt,
		&stats.UpdatedAt,
		&stats.TotalProducts,
		&stats.TotalOrders,
		&stats.TotalRevenue,
		&stats.AverageRating,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}

// SlugExists checks if a slug already exists
func (r *ShopRepository) SlugExists(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM shops WHERE slug = $1 AND ($2::uuid IS NULL OR id != $2))`
	var exists bool
	err := r.db.QueryRow(ctx, query, slug, excludeID).Scan(&exists)
	return exists, err
}
