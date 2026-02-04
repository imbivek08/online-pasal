package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type ReviewHandler struct {
	reviewService *service.ReivewService
	userService   *service.UserService
}

func NewReviewHandler(reviewService *service.ReivewService, userService *service.UserService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		userService:   userService,
	}
}

// CreateReview creates a new product review
func (h *ReviewHandler) CreateReview(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID := middleware.GetClerkUserID(c)
	if clerkID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get user from database
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Parse request body
	var req model.CreateReviewRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Create review
	review, err := h.reviewService.CreateReview(c.Request().Context(), user.ID, &req)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusCreated, "review created successfully", review)
}

// GetReviewByID retrieves a single review by ID
func (h *ReviewHandler) GetReviewByID(c echo.Context) error {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid review ID")
	}

	review, err := h.reviewService.GetReviewByID(c.Request().Context(), reviewID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "review not found")
	}

	return SendSuccess(c, http.StatusOK, "review retrieved successfully", review)
}

// GetProductReviews retrieves all reviews for a product
func (h *ReviewHandler) GetProductReviews(c echo.Context) error {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	// Get pagination parameters
	page := 1
	limit := 10
	if p := c.QueryParam("page"); p != "" {
		if parsedPage, parseErr := strconv.Atoi(p); parseErr == nil && parsedPage > 0 {
			page = parsedPage
		}
	}
	if l := c.QueryParam("limit"); l != "" {
		if parsedLimit, parseErr := strconv.Atoi(l); parseErr == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	reviews, err := h.reviewService.GetProductReviews(c.Request().Context(), productID, page, limit)
	if err != nil {
		// Log the error for debugging
		fmt.Printf("ERROR in GetProductReviews: %v\n", err)
		return SendError(c, http.StatusInternalServerError, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "reviews retrieved successfully", reviews)
}

// GetProductRatingStats retrieves rating statistics for a product
func (h *ReviewHandler) GetProductRatingStats(c echo.Context) error {
	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	stats, err := h.reviewService.GetProductRatingStats(c.Request().Context(), productID)
	if err != nil {
		fmt.Printf("ERROR in GetProductRatingStats: %v\n", err)
		return SendError(c, http.StatusInternalServerError, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "rating stats retrieved successfully", stats)
}

// UpdateReview updates an existing review
func (h *ReviewHandler) UpdateReview(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID := middleware.GetClerkUserID(c)
	if clerkID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get user from database
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid review ID")
	}

	// Parse request body
	var req model.UpdateReviewRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Update review
	err = h.reviewService.UpdateReview(c.Request().Context(), user.ID, reviewID, &req)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "review updated successfully", nil)
}

// DeleteReview deletes a review
func (h *ReviewHandler) DeleteReview(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID := middleware.GetClerkUserID(c)
	if clerkID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get user from database
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid review ID")
	}

	// Delete review
	err = h.reviewService.DeleteReview(c.Request().Context(), user.ID, reviewID)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "review deleted successfully", nil)
}

// MarkReviewHelpful marks a review as helpful
func (h *ReviewHandler) MarkReviewHelpful(c echo.Context) error {
	reviewID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid review ID")
	}

	err = h.reviewService.MarkReviewHelpful(c.Request().Context(), reviewID)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "review marked as helpful", nil)
}

// CanUserReviewProduct checks if a user can review a product
func (h *ReviewHandler) CanUserReviewProduct(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID := middleware.GetClerkUserID(c)
	if clerkID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get user from database
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	productID, err := uuid.Parse(c.Param("productId"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	response, err := h.reviewService.CanUserReviewProduct(c.Request().Context(), user.ID, productID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "check completed", response)
}
