package model

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ShopID        uuid.UUID  `json:"shop_id" db:"shop_id"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	Price         float64    `json:"price" db:"price"`
	StockQuantity int        `json:"stock_quantity" db:"stock_quantity"`
	ImageURL      *string    `json:"image_url,omitempty" db:"image_url"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateProductRequest struct {
	Name          string     `json:"name" validate:"required"`
	Description   *string    `json:"description,omitempty"`
	Price         float64    `json:"price" validate:"required,gt=0"`
	StockQuantity int        `json:"stock_quantity" validate:"gte=0"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	ImageURL      *string    `json:"image_url,omitempty"`
}

type UpdateProductRequest struct {
	Name          *string    `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	Price         *float64   `json:"price,omitempty" validate:"omitempty,gt=0"`
	StockQuantity *int       `json:"stock_quantity,omitempty" validate:"omitempty,gte=0"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	ImageURL      *string    `json:"image_url,omitempty"`
	IsActive      *bool      `json:"is_active,omitempty"`
}

type ProductResponse struct {
	ID            uuid.UUID  `json:"id"`
	ShopID        uuid.UUID  `json:"shop_id"`
	CategoryID    *uuid.UUID `json:"category_id,omitempty"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	Price         float64    `json:"price"`
	StockQuantity int        `json:"stock_quantity"`
	ImageURL      *string    `json:"image_url,omitempty"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (p *Product) ToResponse() *ProductResponse {
	return &ProductResponse{
		ID:            p.ID,
		ShopID:        p.ShopID,
		CategoryID:    p.CategoryID,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		ImageURL:      p.ImageURL,
		IsActive:      p.IsActive,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
