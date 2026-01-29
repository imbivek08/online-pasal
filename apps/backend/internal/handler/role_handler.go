package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
)

type RoleHandler struct {
	userService *service.UserService
}

func NewRoleHandler(userService *service.UserService) *RoleHandler {
	return &RoleHandler{
		userService: userService,
	}
}

// BecomeVendor upgrades a customer to vendor with business information
// POST /api/v1/users/become-vendor
func (h *RoleHandler) BecomeVendor(c echo.Context) error {
	// Get authenticated user from context
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	// Check if already a vendor
	if user.Role == model.RoleVendor {
		return SendError(c, http.StatusBadRequest, nil, "you are already a vendor")
	}

	// Admin cannot become vendor
	if user.Role == model.RoleAdmin {
		return SendError(c, http.StatusBadRequest, nil, "admins cannot become vendors")
	}

	// Parse request body
	var req model.BecomeVendorRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Validate request
	if err := c.Validate(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	// Convert user to vendor
	updatedUser, err := h.userService.ConvertToVendor(c.Request().Context(), user.ID, &req)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to upgrade to vendor")
	}

	// Build response
	response := model.BecomeVendorResponse{
		ID:            updatedUser.ID,
		Email:         updatedUser.Email,
		Role:          updatedUser.Role,
		Phone:         updatedUser.Phone,
		Message:       "Congratulations! You are now a vendor. Create your shop to start selling.",
		CanCreateShop: true,
		NextStep:      "create_shop",
	}

	return SendSuccess(c, http.StatusOK, "successfully upgraded to vendor", response)
}

// GetMyRole returns the current user's role and capabilities
// GET /api/v1/users/my-role
func (h *RoleHandler) GetMyRole(c echo.Context) error {
	user, ok := c.Get("user").(*model.User)
	if !ok || user == nil {
		return SendError(c, http.StatusUnauthorized, nil, "user not found in context")
	}

	roleInfo := map[string]interface{}{
		"role":     user.Role,
		"can_sell": user.Role == model.RoleVendor || user.Role == model.RoleAdmin,
		"can_buy":  true, // Everyone can buy
		"is_admin": user.Role == model.RoleAdmin,
	}

	return SendSuccess(c, http.StatusOK, "role retrieved successfully", roleInfo)
}
