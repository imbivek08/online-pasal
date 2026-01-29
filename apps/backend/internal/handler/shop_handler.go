package handler

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
)

type ShopHandler struct {
	shopService *service.ShopService
}

func NewShopHandler(shopService *service.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

// CreateShop creates a new shop for the authenticated vendor
// POST /api/v1/shops
func (h *ShopHandler) CreateShop(c echo.Context) error {
	// Get vendor ID from context
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "unauthorized")
	}

	// Get user info from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Parse request body
	var req model.CreateShopRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Create shop
	shop, err := h.shopService.CreateShop(c.Request().Context(), user.ID, &req)
	if err != nil {
		if err.Error() == "vendor already has a shop" || err.Error() == "user is not a vendor" {
			return SendError(c, http.StatusBadRequest, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to create shop")
	}

	return SendSuccess(c, http.StatusCreated, "shop created successfully", shop.ToResponse())
}

// GetMyShop retrieves the authenticated vendor's shop
// GET /api/v1/shops/my
func (h *ShopHandler) GetMyShop(c echo.Context) error {
	// Get user info from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Get shop
	shop, err := h.shopService.GetMyShop(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get shop")
	}

	// If no shop found, return null
	if shop == nil {
		return SendSuccess(c, http.StatusOK, "no shop found", nil)
	}

	return SendSuccess(c, http.StatusOK, "shop retrieved successfully", shop.ToResponse())
}

// GetShopByID retrieves a shop by ID
// GET /api/v1/shops/:id
func (h *ShopHandler) GetShopByID(c echo.Context) error {
	// Parse shop ID
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid shop ID")
	}

	// Get shop
	shop, err := h.shopService.GetShopByID(c.Request().Context(), shopID)
	if err != nil {
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to get shop")
	}

	return SendSuccess(c, http.StatusOK, "shop retrieved successfully", shop.ToResponse())
}

// GetShopBySlug retrieves a shop by slug
// GET /api/v1/shops/slug/:slug
func (h *ShopHandler) GetShopBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return SendError(c, http.StatusBadRequest, nil, "slug is required")
	}

	// Get shop
	shop, err := h.shopService.GetShopBySlug(c.Request().Context(), slug)
	if err != nil {
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to get shop")
	}

	return SendSuccess(c, http.StatusOK, "shop retrieved successfully", shop.ToResponse())
}

// ListShops retrieves all shops with pagination
// GET /api/v1/shops
func (h *ShopHandler) ListShops(c echo.Context) error {
	// Parse query parameters
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 20
	}

	search := c.QueryParam("search")
	activeOnly := c.QueryParam("active") == "true"

	// Get shops
	response, err := h.shopService.ListShops(c.Request().Context(), page, pageSize, search, activeOnly)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to list shops")
	}

	return SendSuccess(c, http.StatusOK, "shops retrieved successfully", response)
}

// UpdateShop updates a shop
// PUT /api/v1/shops/:id
func (h *ShopHandler) UpdateShop(c echo.Context) error {
	// Parse shop ID
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid shop ID")
	}

	// Get user from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Parse request body
	var req model.UpdateShopRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Update shop
	shop, err := h.shopService.UpdateShop(c.Request().Context(), shopID, user.ID, &req)
	if err != nil {
		if err.Error() == "unauthorized: you don't own this shop" {
			return SendError(c, http.StatusForbidden, err, "")
		}
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to update shop")
	}

	return SendSuccess(c, http.StatusOK, "shop updated successfully", shop.ToResponse())
}

// ToggleShopStatus toggles the shop's active status
// PATCH /api/v1/shops/:id/status
func (h *ShopHandler) ToggleShopStatus(c echo.Context) error {
	// Parse shop ID
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid shop ID")
	}

	// Get user from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Toggle status
	shop, err := h.shopService.ToggleShopStatus(c.Request().Context(), shopID, user.ID)
	if err != nil {
		if err.Error() == "unauthorized: you don't own this shop" {
			return SendError(c, http.StatusForbidden, err, "")
		}
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to toggle shop status")
	}

	return SendSuccess(c, http.StatusOK, "shop status updated successfully", shop.ToResponse())
}

// GetMyShopStats retrieves statistics for the authenticated vendor's shop
// GET /api/v1/my-shop/stats
func (h *ShopHandler) GetMyShopStats(c echo.Context) error {
	// Get user from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Get shop first
	shop, err := h.shopService.GetMyShop(c.Request().Context(), user.ID)
	if err != nil {
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, nil, "you don't have a shop yet")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to get shop")
	}

	// Get stats
	stats, err := h.shopService.GetShopStats(c.Request().Context(), shop.ID, user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get shop stats")
	}

	return SendSuccess(c, http.StatusOK, "shop stats retrieved successfully", stats)
}

// DeleteShop soft deletes a shop (admin only)
// DELETE /api/v1/shops/:id
func (h *ShopHandler) DeleteShop(c echo.Context) error {
	// Parse shop ID
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid shop ID")
	}

	// Delete shop
	if err := h.shopService.DeleteShop(c.Request().Context(), shopID); err != nil {
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to delete shop")
	}

	return SendSuccess(c, http.StatusOK, "shop deleted successfully", nil)
}

// VerifyShop verifies a shop (admin only)
// PATCH /api/v1/shops/:id/verify
func (h *ShopHandler) VerifyShop(c echo.Context) error {
	// Parse shop ID
	shopIDStr := c.Param("id")
	shopID, err := uuid.Parse(shopIDStr)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid shop ID")
	}

	// Parse request body
	var req struct {
		Verified bool `json:"verified"`
	}
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Verify shop
	shop, err := h.shopService.VerifyShop(c.Request().Context(), shopID, req.Verified)
	if err != nil {
		if err.Error() == "shop not found" {
			return SendError(c, http.StatusNotFound, err, "")
		}
		return SendError(c, http.StatusInternalServerError, err, "failed to verify shop")
	}

	return SendSuccess(c, http.StatusOK, "shop verification updated successfully", shop.ToResponse())
}
