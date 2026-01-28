package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user in the database
func (r *UserRepository) Create(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	user := &model.User{
		ID:        uuid.New(),
		ClerkID:   req.ClerkID,
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		AvatarURL: req.AvatarURL,
		IsActive:  true,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set default role if not provided
	if user.Role == "" {
		user.Role = model.RoleCustomer
	}

	query := `
		INSERT INTO users (id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at, last_login_at
	`

	err := r.db.Pool.QueryRow(ctx, query,
		user.ID, user.ClerkID, user.Email, user.Username, user.FirstName,
		user.LastName, user.Phone, user.AvatarURL, user.IsActive, user.Role,
		user.CreatedAt, user.UpdatedAt,
	).Scan(
		&user.ID, &user.ClerkID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.AvatarURL, &user.IsActive, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByClerkID retrieves a user by their Clerk ID
func (r *UserRepository) GetByClerkID(ctx context.Context, clerkID string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at, last_login_at
		FROM users
		WHERE clerk_id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, clerkID).Scan(
		&user.ID, &user.ClerkID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.AvatarURL, &user.IsActive, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByID retrieves a user by their UUID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.ClerkID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.AvatarURL, &user.IsActive, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user's information
func (r *UserRepository) Update(ctx context.Context, clerkID string, req *model.UpdateUserRequest) (*model.User, error) {
	query := `
		UPDATE users
		SET username = COALESCE($2, username),
		    first_name = COALESCE($3, first_name),
		    last_name = COALESCE($4, last_name),
		    phone = COALESCE($5, phone),
		    avatar_url = COALESCE($6, avatar_url),
		    role = COALESCE($7, role),
		    updated_at = $8
		WHERE clerk_id = $1
		RETURNING id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at, last_login_at
	`

	user := &model.User{}
	err := r.db.Pool.QueryRow(ctx, query,
		clerkID, req.Username, req.FirstName, req.LastName, req.Phone, req.AvatarURL, req.Role, time.Now(),
	).Scan(
		&user.ID, &user.ClerkID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.AvatarURL, &user.IsActive, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// UpdateLastLogin updates the user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, clerkID string) error {
	query := `
		UPDATE users
		SET last_login_at = $2
		WHERE clerk_id = $1
	`

	_, err := r.db.Pool.Exec(ctx, query, clerkID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}

// Delete soft deletes a user by setting is_active to false
func (r *UserRepository) Delete(ctx context.Context, clerkID string) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = $2
		WHERE clerk_id = $1
	`

	result, err := r.db.Pool.Exec(ctx, query, clerkID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	query := `
		SELECT id, clerk_id, email, username, first_name, last_name, phone, avatar_url, is_active, role, created_at, updated_at, last_login_at
		FROM users
		WHERE email = $1
	`

	err := r.db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.ClerkID, &user.Email, &user.Username, &user.FirstName,
		&user.LastName, &user.Phone, &user.AvatarURL, &user.IsActive, &user.Role,
		&user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetShopIDByVendorID retrieves a shop ID for a vendor
func (r *UserRepository) GetShopIDByVendorID(ctx context.Context, vendorID uuid.UUID) (uuid.UUID, error) {
	var shopID uuid.UUID

	query := `
		SELECT id
		FROM shops
		WHERE vendor_id = $1 AND is_active = true
		LIMIT 1
	`

	err := r.db.Pool.QueryRow(ctx, query, vendorID).Scan(&shopID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return uuid.Nil, fmt.Errorf("no shop found for vendor")
		}
		return uuid.Nil, fmt.Errorf("failed to get shop ID: %w", err)
	}

	return shopID, nil
}
