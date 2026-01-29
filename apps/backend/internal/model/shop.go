package model

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ID          uuid.UUID `json:"id"`
	VendorID    uuid.UUID `json:"vendor_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	LogoURL     *string   `json:"logo_url"`
	BannerURL   *string   `json:"banner_url"`
	Address     *string   `json:"address"`
	City        *string   `json:"city"`
	State       *string   `json:"state"`
	Country     *string   `json:"country"`
	PostalCode  *string   `json:"postal_code"`
	Phone       *string   `json:"phone"`
	Email       *string   `json:"email"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateShopRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Address     *string `json:"address" validate:"omitempty,max=255"`
	City        *string `json:"city" validate:"omitempty,max=100"`
	State       *string `json:"state" validate:"omitempty,max=100"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=20"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	Email       *string `json:"email" validate:"omitempty,email"`
}

type UpdateShopRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=3,max=100"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	LogoURL     *string `json:"logo_url" validate:"omitempty,url"`
	BannerURL   *string `json:"banner_url" validate:"omitempty,url"`
	Address     *string `json:"address" validate:"omitempty,max=255"`
	City        *string `json:"city" validate:"omitempty,max=100"`
	State       *string `json:"state" validate:"omitempty,max=100"`
	Country     *string `json:"country" validate:"omitempty,max=100"`
	PostalCode  *string `json:"postal_code" validate:"omitempty,max=20"`
	Phone       *string `json:"phone" validate:"omitempty,max=20"`
	Email       *string `json:"email" validate:"omitempty,email"`
}

type ShopResponse struct {
	ID          uuid.UUID `json:"id"`
	VendorID    uuid.UUID `json:"vendor_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	LogoURL     *string   `json:"logo_url"`
	BannerURL   *string   `json:"banner_url"`
	Address     *string   `json:"address"`
	City        *string   `json:"city"`
	State       *string   `json:"state"`
	Country     *string   `json:"country"`
	PostalCode  *string   `json:"postal_code"`
	Phone       *string   `json:"phone"`
	Email       *string   `json:"email"`
	IsActive    bool      `json:"is_active"`
	IsVerified  bool      `json:"is_verified"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShopWithStats struct {
	ShopResponse
	TotalProducts int     `json:"total_products"`
	TotalOrders   int     `json:"total_orders"`
	TotalRevenue  float64 `json:"total_revenue"`
	AverageRating float64 `json:"average_rating"`
}

type ShopListResponse struct {
	Shops      []ShopResponse `json:"shops"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// ToResponse converts Shop to ShopResponse
func (s *Shop) ToResponse() ShopResponse {
	return ShopResponse{
		ID:          s.ID,
		VendorID:    s.VendorID,
		Name:        s.Name,
		Slug:        s.Slug,
		Description: s.Description,
		LogoURL:     s.LogoURL,
		BannerURL:   s.BannerURL,
		Address:     s.Address,
		City:        s.City,
		State:       s.State,
		Country:     s.Country,
		PostalCode:  s.PostalCode,
		Phone:       s.Phone,
		Email:       s.Email,
		IsActive:    s.IsActive,
		IsVerified:  s.IsVerified,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}
