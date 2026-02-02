package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	orderService *service.OrderService
	userService  *service.UserService
	shopService  *service.ShopService
}

func NewOrderHandler(orderService *service.OrderService, userService *service.UserService, shopService *service.ShopService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		userService:  userService,
		shopService:  shopService,
	}
}

// CreateOrder creates an order from cart (checkout)
func (h *OrderHandler) CreateOrder(c echo.Context) error {
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
	var req model.CreateOrderRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Create order
	order, err := h.orderService.CreateOrderFromCart(c.Request().Context(), user.ID, &req)
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusCreated, "order created successfully", order)
}

// GetOrders retrieves user's order history
func (h *OrderHandler) GetOrders(c echo.Context) error {
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

	// Get orders
	orders, err := h.orderService.GetUserOrders(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get orders")
	}

	return SendSuccess(c, http.StatusOK, "orders retrieved successfully", orders)
}

// GetOrderByID retrieves a single order by ID
func (h *OrderHandler) GetOrderByID(c echo.Context) error {
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

	// Parse order ID
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid order ID")
	}

	// Get order
	order, err := h.orderService.GetOrderByID(c.Request().Context(), orderID, user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "order not found")
	}

	return SendSuccess(c, http.StatusOK, "order retrieved successfully", order)
}

// GetVendorOrders retrieves orders for vendor's shop
func (h *OrderHandler) GetVendorOrders(c echo.Context) error {
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

	// Check if user is vendor
	if user.Role != model.RoleVendor {
		return SendError(c, http.StatusForbidden, nil, "vendor access required")
	}

	// Get shop ID for vendor
	shopID, err := h.userService.GetShopIDByVendorID(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "no shop found for vendor")
	}

	// Get orders
	orders, err := h.orderService.GetVendorOrders(c.Request().Context(), shopID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get orders")
	}

	return SendSuccess(c, http.StatusOK, "vendor orders retrieved successfully", orders)
}

// UpdateOrderStatus updates order status (vendor only)
func (h *OrderHandler) UpdateOrderStatus(c echo.Context) error {
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

	// Check if user is vendor
	if user.Role != model.RoleVendor {
		return SendError(c, http.StatusForbidden, nil, "vendor access required")
	}

	// Parse order ID
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid order ID")
	}

	// Parse request
	var req model.UpdateOrderStatusRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Update order status
	if err := h.orderService.UpdateOrderStatus(c.Request().Context(), orderID, req.Status); err != nil {
		return SendError(c, http.StatusBadRequest, err, "failed to update order status")
	}

	return SendSuccess(c, http.StatusOK, "order status updated successfully", nil)
}

// CancelOrder cancels an order
func (h *OrderHandler) CancelOrder(c echo.Context) error {
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

	// Parse order ID
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid order ID")
	}

	// Cancel order
	if err := h.orderService.CancelOrder(c.Request().Context(), orderID, user.ID); err != nil {
		return SendError(c, http.StatusBadRequest, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "order cancelled successfully", nil)
}
