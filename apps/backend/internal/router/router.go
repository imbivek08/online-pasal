package router

import (
	"net/http"

	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/handler"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/repository"
	"github.com/imbivek08/hamropasal/internal/service"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *database.Database, cfg *config.Config) {
	// Health check endpoint
	e.GET("/health", healthCheck)

	// Initialize Clerk client
	middleware.InitClerkClient(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService, userService)
	webhookHandler := handler.NewWebhookHandler(userService)

	// API v1 group
	v1 := e.Group("/api/v1")

	// Webhook routes (no auth required)
	webhooks := v1.Group("/webhooks")
	webhooks.POST("/clerk", webhookHandler.HandleClerkWebhook)

	// Auth middleware for protected routes
	authMiddleware := middleware.ClerkAuthMiddleware(cfg)

	// User routes (protected)
	users := v1.Group("/users", authMiddleware)
	users.GET("/profile", userHandler.GetProfile)
	users.PUT("/profile", userHandler.UpdateProfile)
	users.DELETE("/account", userHandler.DeleteAccount)
	users.GET("/:id", userHandler.GetUserByID)

	// Product routes
	setupProductRoutes(v1, productHandler, authMiddleware)

	// Shop routes (to be implemented)
	// setupShopRoutes(v1, db, authMiddleware)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "server is running",
	})
}

// Placeholder route groups - implement these as you add features
func setupUserRoutes(g *echo.Group, db *database.Database) {
	// Moved to main SetupRoutes function
}

func setupProductRoutes(g *echo.Group, productHandler *handler.ProductHandler, authMiddleware echo.MiddlewareFunc) {
	products := g.Group("/products")

	// Public routes
	products.GET("", productHandler.GetProducts)        // Get all products
	products.GET("/:id", productHandler.GetProductByID) // Get single product

	// Protected routes (vendor only)
	products.POST("", productHandler.CreateProduct, authMiddleware)       // Create product
	products.PUT("/:id", productHandler.UpdateProduct, authMiddleware)    // Update product
	products.DELETE("/:id", productHandler.DeleteProduct, authMiddleware) // Delete product

	// Vendor-specific routes
	vendor := g.Group("/vendor", authMiddleware)
	vendor.GET("/products", productHandler.GetVendorProducts) // Get my products
}

// func setupShopRoutes(g *echo.Group, db *database.Database, authMiddleware echo.MiddlewareFunc) {
// 	// shops := g.Group("/shops")
// 	// Add shop-related routes here
// }
