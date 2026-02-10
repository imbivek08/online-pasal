package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type AddressHandler struct {
	addressService *service.AddressService
	userService    *service.UserService
}

func NewAddressHandler(addressService *service.AddressService, userService *service.UserService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
		userService:    userService,
	}
}

// GetAddresses returns all addresses for the authenticated user
func (h *AddressHandler) GetAddresses(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	addresses, err := h.addressService.GetUserAddresses(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to get addresses")
	}

	return SendSuccess(c, http.StatusOK, "addresses retrieved", addresses)
}

// GetAddress returns a single address by ID
func (h *AddressHandler) GetAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid address ID")
	}

	address, err := h.addressService.GetAddressByID(c.Request().Context(), addressID, user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "address not found")
	}

	return SendSuccess(c, http.StatusOK, "address retrieved", address)
}

// GetDefaultAddress returns the user's default address
func (h *AddressHandler) GetDefaultAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	address, err := h.addressService.GetDefaultAddress(c.Request().Context(), user.ID)
	if err != nil {
		return SendSuccess(c, http.StatusOK, "no default address found", nil)
	}

	return SendSuccess(c, http.StatusOK, "default address retrieved", address)
}

// CreateAddress creates a new address
func (h *AddressHandler) CreateAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	var input model.AddressInput
	if err := c.Bind(&input); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&input); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	address, err := h.addressService.CreateAddress(c.Request().Context(), user.ID, &input)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to create address")
	}

	return SendSuccess(c, http.StatusCreated, "address created", address)
}

// UpdateAddress updates an existing address
func (h *AddressHandler) UpdateAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid address ID")
	}

	var input model.AddressInput
	if err := c.Bind(&input); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	if err := c.Validate(&input); err != nil {
		return SendError(c, http.StatusBadRequest, err, "validation failed")
	}

	address, err := h.addressService.UpdateAddress(c.Request().Context(), addressID, user.ID, &input)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "address updated", address)
}

// DeleteAddress deletes an address
func (h *AddressHandler) DeleteAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid address ID")
	}

	if err := h.addressService.DeleteAddress(c.Request().Context(), addressID, user.ID); err != nil {
		return SendError(c, http.StatusNotFound, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "address deleted", nil)
}

// SetDefaultAddress sets an address as the default
func (h *AddressHandler) SetDefaultAddress(c echo.Context) error {
	clerkUserID := middleware.GetClerkUserID(c)
	if clerkUserID == "" {
		return SendError(c, http.StatusUnauthorized, nil, "user not authenticated")
	}

	user, err := h.userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "user not found")
	}

	addressID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid address ID")
	}

	if err := h.addressService.SetDefaultAddress(c.Request().Context(), addressID, user.ID); err != nil {
		return SendError(c, http.StatusNotFound, err, err.Error())
	}

	return SendSuccess(c, http.StatusOK, "default address updated", nil)
}
