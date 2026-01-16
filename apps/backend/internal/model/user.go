package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleCustomer UserRole = "customer"
	RoleVendor   UserRole = "vendor"
	RoleAdmin    UserRole = "admin"
)

type User struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ClerkID     string     `json:"clerk_id" db:"clerk_id"`
	Email       string     `json:"email" db:"email"`
	Username    *string    `json:"username,omitempty" db:"username"`
	FirstName   *string    `json:"first_name,omitempty" db:"first_name"`
	LastName    *string    `json:"last_name,omitempty" db:"last_name"`
	Phone       *string    `json:"phone,omitempty" db:"phone"`
	AvatarURL   *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	Role        UserRole   `json:"role" db:"role"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty" db:"last_login_at"`
}

type CreateUserRequest struct {
	ClerkID   string   `json:"clerk_id" validate:"required"`
	Email     string   `json:"email" validate:"required,email"`
	Username  *string  `json:"username,omitempty"`
	FirstName *string  `json:"first_name,omitempty"`
	LastName  *string  `json:"last_name,omitempty"`
	Phone     *string  `json:"phone,omitempty"`
	AvatarURL *string  `json:"avatar_url,omitempty"`
	Role      UserRole `json:"role,omitempty"`
}

type UpdateUserRequest struct {
	Username  *string  `json:"username,omitempty"`
	FirstName *string  `json:"first_name,omitempty"`
	LastName  *string  `json:"last_name,omitempty"`
	Phone     *string  `json:"phone,omitempty"`
	AvatarURL *string  `json:"avatar_url,omitempty"`
	Role      UserRole `json:"role,omitempty"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  *string   `json:"username,omitempty"`
	FirstName *string   `json:"first_name,omitempty"`
	LastName  *string   `json:"last_name,omitempty"`
	Phone     *string   `json:"phone,omitempty"`
	AvatarURL *string   `json:"avatar_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		AvatarURL: u.AvatarURL,
		IsActive:  u.IsActive,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
