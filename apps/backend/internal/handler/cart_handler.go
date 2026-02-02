package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	cartService *service.CartService
	userService *service.UserService
}

func NewCartHandler(cartService *service.CartService, userService *service.UserService) *CartHandler {
	return &CartHandler{
		cartService: cartService,
		userService: userService,
	}
}

// GetCart retrieves user's cart with all items
func (h *CartHandler) GetCart(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Get cart
	cart, err := h.cartService.GetCart(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get cart")
	}

	return SendSuccess(c, http.StatusOK, "cart retrieved successfully", cart)
}

// AddToCart adds a product to cart
func (h *CartHandler) AddToCart(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Parse request
	var req model.AddToCartRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Add to cart
	item, err := h.cartService.AddToCart(c.Request().Context(), user.ID, &req)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusCreated, "product added to cart", item)
}

// UpdateCartItem updates cart item quantity
func (h *CartHandler) UpdateCartItem(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Parse item ID
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid item ID")
	}

	// Parse request
	var req model.UpdateCartItemRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Update item quantity
	if err := h.cartService.UpdateCartItemQuantity(c.Request().Context(), user.ID, itemID, req.Quantity); err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "cart item updated successfully", nil)
}

// RemoveCartItem removes an item from cart
func (h *CartHandler) RemoveCartItem(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Parse item ID
	itemID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid item ID")
	}

	// Remove item
	if err := h.cartService.RemoveCartItem(c.Request().Context(), user.ID, itemID); err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "item removed from cart successfully", nil)
}

// ClearCart removes all items from cart
func (h *CartHandler) ClearCart(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Clear cart
	if err := h.cartService.ClearCart(c.Request().Context(), user.ID); err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to clear cart")
	}

	return SendSuccess(c, http.StatusOK, "cart cleared successfully", nil)
}

// GetCartItemCount returns the total number of items in cart
func (h *CartHandler) GetCartItemCount(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get internal user
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	// Get item count
	count, err := h.cartService.GetCartItemCount(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get cart item count")
	}

	return SendSuccess(c, http.StatusOK, "cart item count retrieved successfully", map[string]int{"count": count})
}
