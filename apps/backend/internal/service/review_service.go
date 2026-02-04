package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/repository"
)

type ReivewService struct {
	reviewRepo  *repository.ReviewRepository
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
}

func NewReviewService(
	reviewRepo *repository.ReviewRepository,
	orderRepo *repository.OrderRepository,
	productRepo *repository.ProductRepository,
) *ReivewService {
	return &ReivewService{
		reviewRepo:  reviewRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

func (s *ReivewService) CreateReview(ctx context.Context, userID uuid.UUID, req *model.CreateReviewRequest) (*model.Review, error) {
	// 1. Verify the product exists
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if !product.IsActive {
		return nil, errors.New("product is not active")
	}

	// 2. Verify the order exists and belongs to the user
	order, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.UserID != userID {
		return nil, errors.New("unauthorized: order does not belong to user")
	}

	// 3. Verify the order is delivered (can only review delivered orders)
	if order.Status != model.OrderStatusDelivered {
		return nil, errors.New("can only review delivered orders")
	}

	// 4. Verify the order contains the product being reviewed
	orderItems, err := s.orderRepo.GetOrderItems(ctx, req.OrderID)
	if err != nil {
		return nil, errors.New("failed to get order items")
	}

	productInOrder := false
	for _, item := range orderItems {
		if item.ProductID == req.ProductID {
			productInOrder = true
			break
		}
	}
	if !productInOrder {
		return nil, errors.New("product not found in order")
	}

	// 5. Create the review
	now := time.Now()
	review := &model.Review{
		ID:                 uuid.New(),
		ProductID:          req.ProductID,
		UserID:             userID,
		OrderID:            &req.OrderID,
		Rating:             req.Rating,
		Title:              req.Title,
		Comment:            req.Comment,
		IsVerifiedPurchase: true, // Since we verified the order
		IsApproved:         true, // Auto-approve for now (can add moderation later)
		HelpfulCount:       0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err = s.reviewRepo.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}

// GetReviewByID retrieves a review by ID
func (s *ReivewService) GetReviewByID(ctx context.Context, id uuid.UUID) (*model.ReviewWithUser, error) {
	return s.reviewRepo.GetByID(ctx, id)
}

// GetProductReviews retrieves reviews for a product with pagination
func (s *ReivewService) GetProductReviews(ctx context.Context, productID uuid.UUID, page, limit int) (*model.ReviewListResponse, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Calculate offset
	offset := (page - 1) * limit

	// Get reviews
	reviews, err := s.reviewRepo.GetByProductID(ctx, productID, limit, offset)
	if err != nil {
		return nil, err
	}

	// Get total count
	totalReviews, err := s.reviewRepo.CountByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return &model.ReviewListResponse{
		Reviews:      reviews,
		TotalReviews: totalReviews,
		Page:         page,
		Limit:        limit,
	}, nil
}

// GetProductRatingStats retrieves rating statistics for a product
func (s *ReivewService) GetProductRatingStats(ctx context.Context, productID uuid.UUID) (*model.ProductRatingStats, error) {
	// Verify product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	return s.reviewRepo.GetProductRatingStats(ctx, productID)
}

// UpdateReview updates an existing review
func (s *ReivewService) UpdateReview(ctx context.Context, userID, reviewID uuid.UUID, req *model.UpdateReviewRequest) error {
	// Get the review to verify ownership
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("unauthorized: you can only update your own reviews")
	}

	return s.reviewRepo.UpdateReview(ctx, reviewID, req)
}

// DeleteReview deletes a review
func (s *ReivewService) DeleteReview(ctx context.Context, userID, reviewID uuid.UUID) error {
	// Get the review to verify ownership
	review, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("unauthorized: you can only delete your own reviews")
	}

	return s.reviewRepo.DeleteReview(ctx, reviewID)
}

// MarkReviewHelpful increments the helpful count for a review
func (s *ReivewService) MarkReviewHelpful(ctx context.Context, reviewID uuid.UUID) error {
	// Verify review exists
	_, err := s.reviewRepo.GetByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	return s.reviewRepo.IncrementHelpfulCount(ctx, reviewID)
}

// CanUserReviewProduct checks if a user can review a product
func (s *ReivewService) CanUserReviewProduct(ctx context.Context, userID, productID uuid.UUID) (*model.CanReviewResponse, error) {
	response := &model.CanReviewResponse{
		CanReview: false,
		Reason:    "",
	}

	// Check if product exists
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		response.Reason = "Product not found"
		return response, nil
	}

	// Check if user already reviewed this product
	hasReviewed, err := s.reviewRepo.HasUserReviewedProduct(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	if hasReviewed {
		response.Reason = "You have already reviewed this product"
		return response, nil
	}

	// Check if user has a delivered order with this product
	// This requires getting user's orders and checking
	// For simplicity, we'll allow if they haven't reviewed yet
	response.CanReview = true
	return response, nil
}
