package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	productService *service.ProductService
	userService    *service.UserService
}

func NewProductHandler(productService *service.ProductService, userService *service.UserService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		userService:    userService,
	}
}

// GetProducts retrieves all products (public)
func (h *ProductHandler) GetProducts(c echo.Context) error {
	products, err := h.productService.GetAllProducts(c.Request().Context(), nil)
	fmt.Println("this end point was hitted")
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to retrieve products")
	}

	// Convert to response format
	var responses []*model.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return SendSuccess(c, http.StatusOK, "products retrieved successfully", responses)
}

// GetProductByID retrieves a single product by ID (public)
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	product, err := h.productService.GetProductByID(c.Request().Context(), productID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "product not found")
	}

	return SendSuccess(c, http.StatusOK, "product retrieved successfully", product.ToResponse())
}

// CreateProduct creates a new product (vendor only)
func (h *ProductHandler) CreateProduct(c echo.Context) error {
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

	// Check if user is a vendor
	if user.Role != model.RoleVendor {
		return SendError(c, http.StatusForbidden, nil, "only vendors can create products")
	}

	// Get shop ID for vendor
	shopID, err := h.userService.GetShopIDByVendorID(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "no shop found for vendor")
	}

	// Parse request body
	var req model.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Create product
	product, err := h.productService.CreateProduct(c.Request().Context(), shopID, &req)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to create product")
	}

	return SendSuccess(c, http.StatusCreated, "product created successfully", product.ToResponse())
}

// UpdateProduct updates a product (vendor only, own products)
func (h *ProductHandler) UpdateProduct(c echo.Context) error {
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

	// Get shop ID for vendor
	shopID, err := h.userService.GetShopIDByVendorID(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "no shop found for vendor")
	}

	// Parse product ID
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	// Parse request body
	var req model.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid request body")
	}

	// Update product
	product, err := h.productService.UpdateProduct(c.Request().Context(), productID, shopID, &req)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to update product")
	}

	return SendSuccess(c, http.StatusOK, "product updated successfully", product.ToResponse())
}

// DeleteProduct deletes a product (vendor only, own products)
func (h *ProductHandler) DeleteProduct(c echo.Context) error {
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

	// Get shop ID for vendor
	shopID, err := h.userService.GetShopIDByVendorID(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "no shop found for vendor")
	}

	// Parse product ID
	productID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return SendError(c, http.StatusBadRequest, err, "invalid product ID")
	}

	// Delete product
	if err := h.productService.DeleteProduct(c.Request().Context(), productID, shopID); err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to delete product")
	}

	return SendSuccess(c, http.StatusOK, "product deleted successfully", nil)
}

// GetVendorProducts retrieves all products for the authenticated vendor
func (h *ProductHandler) GetVendorProducts(c echo.Context) error {
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

	// Get shop ID for vendor
	shopID, err := h.userService.GetShopIDByVendorID(c.Request().Context(), user.ID)
	if err != nil {
		return SendError(c, http.StatusNotFound, err, "no shop found for vendor")
	}

	// Get shop's products
	products, err := h.productService.GetShopProducts(c.Request().Context(), shopID)
	if err != nil {
		return SendError(c, http.StatusInternalServerError, err, "failed to retrieve products")
	}

	// Convert to response format
	var responses []*model.ProductResponse
	for _, product := range products {
		responses = append(responses, product.ToResponse())
	}

	return SendSuccess(c, http.StatusOK, "vendor products retrieved successfully", responses)
}
