package model

import (
	"time"

	"github.com/google/uuid"
)

// Review represents a product review from a user
type Review struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	ProductID          uuid.UUID  `json:"product_id" db:"product_id"`
	UserID             uuid.UUID  `json:"user_id" db:"user_id"`
	OrderID            *uuid.UUID `json:"order_id,omitempty" db:"order_id"`
	Rating             int        `json:"rating" db:"rating"`
	Title              *string    `json:"title,omitempty" db:"title"`
	Comment            *string    `json:"comment,omitempty" db:"comment"`
	IsVerifiedPurchase bool       `json:"is_verified_purchase" db:"is_verified_purchase"`
	IsApproved         bool       `json:"is_approved" db:"is_approved"`
	HelpfulCount       int        `json:"helpful_count" db:"helpful_count"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// ReviewWithUser extends Review with user information
type ReviewWithUser struct {
	Review
	UserName   *string `json:"user_name,omitempty" db:"user_name"`
	UserAvatar *string `json:"user_avatar,omitempty" db:"user_avatar"`
}

// CreateReviewRequest represents the request to create a new review
type CreateReviewRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	OrderID   uuid.UUID `json:"order_id" validate:"required"`
	Rating    int       `json:"rating" validate:"required,min=1,max=5"`
	Title     *string   `json:"title,omitempty" validate:"omitempty,max=200"`
	Comment   *string   `json:"comment,omitempty" validate:"omitempty,max=2000"`
}

// UpdateReviewRequest represents the request to update an existing review
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	Title   *string `json:"title,omitempty" validate:"omitempty,max=200"`
	Comment *string `json:"comment,omitempty" validate:"omitempty,max=2000"`
}

// ProductRatingStats represents aggregated rating statistics for a product
type ProductRatingStats struct {
	ProductID      uuid.UUID `json:"product_id"`
	AverageRating  float64   `json:"average_rating"`
	TotalReviews   int       `json:"total_reviews"`
	FiveStarCount  int       `json:"five_star_count"`
	FourStarCount  int       `json:"four_star_count"`
	ThreeStarCount int       `json:"three_star_count"`
	TwoStarCount   int       `json:"two_star_count"`
	OneStarCount   int       `json:"one_star_count"`
}

// ReviewListResponse represents paginated review list response
type ReviewListResponse struct {
	Reviews      []ReviewWithUser `json:"reviews"`
	TotalReviews int              `json:"total_reviews"`
	Page         int              `json:"page"`
	Limit        int              `json:"limit"`
}

// CanReviewResponse indicates whether a user can review a product
type CanReviewResponse struct {
	CanReview        bool       `json:"can_review"`
	Reason           string     `json:"reason,omitempty"`
	ExistingReviewID *uuid.UUID `json:"existing_review_id,omitempty"`
}
