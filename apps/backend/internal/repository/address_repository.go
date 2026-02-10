package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
)

type AddressRepository struct {
	db *database.Database
}

func NewAddressRepository(db *database.Database) *AddressRepository {
	return &AddressRepository{db: db}
}

// GetByUserID retrieves all addresses for a user
func (r *AddressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*model.Address, error) {
	var addresses []*model.Address
	query := `
		SELECT id, user_id, full_name, phone, address_line1, address_line2,
		       city, state, postal_code, country, is_default, address_type,
		       created_at, updated_at
		FROM addresses
		WHERE user_id = $1
		ORDER BY is_default DESC, created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var addr model.Address
		err := rows.Scan(
			&addr.ID,
			&addr.UserID,
			&addr.FullName,
			&addr.Phone,
			&addr.AddressLine1,
			&addr.AddressLine2,
			&addr.City,
			&addr.State,
			&addr.PostalCode,
			&addr.Country,
			&addr.IsDefault,
			&addr.AddressType,
			&addr.CreatedAt,
			&addr.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, &addr)
	}

	return addresses, rows.Err()
}

// GetByID retrieves an address by ID
func (r *AddressRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Address, error) {
	var addr model.Address
	query := `
		SELECT id, user_id, full_name, phone, address_line1, address_line2,
		       city, state, postal_code, country, is_default, address_type,
		       created_at, updated_at
		FROM addresses
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&addr.ID,
		&addr.UserID,
		&addr.FullName,
		&addr.Phone,
		&addr.AddressLine1,
		&addr.AddressLine2,
		&addr.City,
		&addr.State,
		&addr.PostalCode,
		&addr.Country,
		&addr.IsDefault,
		&addr.AddressType,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// GetDefaultByUserID retrieves the user's default address
func (r *AddressRepository) GetDefaultByUserID(ctx context.Context, userID uuid.UUID) (*model.Address, error) {
	var addr model.Address
	query := `
		SELECT id, user_id, full_name, phone, address_line1, address_line2,
		       city, state, postal_code, country, is_default, address_type,
		       created_at, updated_at
		FROM addresses
		WHERE user_id = $1 AND is_default = true
		LIMIT 1
	`

	err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&addr.ID,
		&addr.UserID,
		&addr.FullName,
		&addr.Phone,
		&addr.AddressLine1,
		&addr.AddressLine2,
		&addr.City,
		&addr.State,
		&addr.PostalCode,
		&addr.Country,
		&addr.IsDefault,
		&addr.AddressType,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &addr, nil
}

// Create creates a new address
func (r *AddressRepository) Create(ctx context.Context, addr *model.Address) error {
	query := `
		INSERT INTO addresses (
			id, user_id, full_name, phone, address_line1, address_line2,
			city, state, postal_code, country, is_default, address_type,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.Pool.Exec(ctx, query,
		addr.ID,
		addr.UserID,
		addr.FullName,
		addr.Phone,
		addr.AddressLine1,
		addr.AddressLine2,
		addr.City,
		addr.State,
		addr.PostalCode,
		addr.Country,
		addr.IsDefault,
		addr.AddressType,
		addr.CreatedAt,
		addr.UpdatedAt,
	)
	return err
}

// Update updates an existing address
func (r *AddressRepository) Update(ctx context.Context, addr *model.Address) error {
	query := `
		UPDATE addresses
		SET full_name = $1, phone = $2, address_line1 = $3, address_line2 = $4,
		    city = $5, state = $6, postal_code = $7, country = $8,
		    address_type = $9, updated_at = $10
		WHERE id = $11 AND user_id = $12
	`

	result, err := r.db.Pool.Exec(ctx, query,
		addr.FullName,
		addr.Phone,
		addr.AddressLine1,
		addr.AddressLine2,
		addr.City,
		addr.State,
		addr.PostalCode,
		addr.Country,
		addr.AddressType,
		addr.UpdatedAt,
		addr.ID,
		addr.UserID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("address not found or not owned by user")
	}
	return nil
}

// Delete deletes an address
func (r *AddressRepository) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM addresses WHERE id = $1 AND user_id = $2`
	result, err := r.db.Pool.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("address not found or not owned by user")
	}
	return nil
}

// UnsetDefaultForUser unsets any default address for a user
func (r *AddressRepository) UnsetDefaultForUser(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE addresses SET is_default = false, updated_at = $1 WHERE user_id = $2 AND is_default = true`
	_, err := r.db.Pool.Exec(ctx, query, time.Now(), userID)
	return err
}

// SetDefault sets an address as the default for a user
func (r *AddressRepository) SetDefault(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	// First unset any existing default
	if err := r.UnsetDefaultForUser(ctx, userID); err != nil {
		return err
	}

	// Set the new default
	query := `UPDATE addresses SET is_default = true, updated_at = $1 WHERE id = $2 AND user_id = $3`
	result, err := r.db.Pool.Exec(ctx, query, time.Now(), id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("address not found or not owned by user")
	}
	return nil
}
