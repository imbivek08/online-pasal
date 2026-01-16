package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile gets the authenticated user's profile
func (h *UserHandler) GetProfile(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID, ok := middleware.GetClerkUserID(c)
	if !ok {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Get user from database
	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	return SendSuccess(c, http.StatusOK, "profile retrieved successfully", user.ToResponse())
}

// UpdateProfile updates the authenticated user's profile
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID, ok := middleware.GetClerkUserID(c)
	if !ok {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Parse request body
	var req model.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return SendValidationError(c, "invalid request body")
	}

	// Update user
	user, err := h.userService.UpdateUser(c.Request().Context(), clerkID, &req)
	if err != nil {
		return SendInternalError(c, err)
	}

	return SendSuccess(c, http.StatusOK, "profile updated successfully", user.ToResponse())
}

// GetUserByID gets a user by their ID (admin or public endpoint)
func (h *UserHandler) GetUserByID(c echo.Context) error {
	// Parse user ID from URL parameter
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return SendValidationError(c, "invalid user ID")
	}

	// Get user from database
	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	return SendSuccess(c, http.StatusOK, "user retrieved successfully", user.ToResponse())
}

// DeleteAccount soft deletes the authenticated user's account
func (h *UserHandler) DeleteAccount(c echo.Context) error {
	// Get Clerk user ID from middleware
	clerkID, ok := middleware.GetClerkUserID(c)
	if !ok {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	// Delete user
	err := h.userService.DeleteUser(c.Request().Context(), clerkID)
	if err != nil {
		return SendInternalError(c, err)
	}

	return SendSuccess(c, http.StatusOK, "account deleted successfully", nil)
}
