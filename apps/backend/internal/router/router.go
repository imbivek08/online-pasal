package router

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/imbivek08/hamropasal/internal/handler"
	"github.com/imbivek08/hamropasal/internal/middleware"
	"github.com/imbivek08/hamropasal/internal/repository"
	"github.com/imbivek08/hamropasal/internal/service"
)

func SetupRoutes(e *echo.Echo, db *database.Database, cfg *config.Config) {
	// Health check endpoint
	e.GET("/health", healthCheck)

	// Initialize Clerk client
	middleware.InitClerkClient(cfg)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	shopRepo := repository.NewShopRepository(db.Pool)
	cartRepo := repository.NewCartRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	productService := service.NewProductService(productRepo)
	shopService := service.NewShopService(shopRepo, userRepo)
	cartService := service.NewCartService(cartRepo, productRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, productRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	productHandler := handler.NewProductHandler(productService, userService)
	shopHandler := handler.NewShopHandler(shopService)
	roleHandler := handler.NewRoleHandler(userService)
	webhookHandler := handler.NewWebhookHandler(userService)
	cartHandler := handler.NewCartHandler(cartService, userService)
	orderHandler := handler.NewOrderHandler(orderService, userService, shopService)

	// API v1 group
	v1 := e.Group("/api/v1")

	// Webhook routes (no auth required)
	webhooks := v1.Group("/webhooks")
	webhooks.POST("/clerk", webhookHandler.HandleClerkWebhook)

	// Auth middleware for protected routes
	authMiddleware := middleware.ClerkAuthMiddleware(cfg)
	loadUserMiddleware := middleware.LoadUserMiddleware(userService)

	// User routes (protected)
	users := v1.Group("/users", authMiddleware, loadUserMiddleware)
	users.GET("/profile", userHandler.GetProfile)
	users.PUT("/profile", userHandler.UpdateProfile)
	users.DELETE("/account", userHandler.DeleteAccount)
	users.GET("/:id", userHandler.GetUserByID)
	users.POST("/become-vendor", roleHandler.BecomeVendor)
	users.GET("/my-role", roleHandler.GetMyRole)

	// Product routes
	setupProductRoutes(v1, productHandler, authMiddleware, loadUserMiddleware)

	// Shop routes
	setupShopRoutes(v1, shopHandler, authMiddleware, loadUserMiddleware)

	// Cart routes
	setupCartRoutes(v1, cartHandler, authMiddleware, loadUserMiddleware)

	// Order routes
	setupOrderRoutes(v1, orderHandler, authMiddleware, loadUserMiddleware)
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

func setupProductRoutes(g *echo.Group, productHandler *handler.ProductHandler, authMiddleware, loadUserMiddleware echo.MiddlewareFunc) {
	products := g.Group("/products")

	// Public routes
	products.GET("", productHandler.GetProducts)        // Get all products
	products.GET("/:id", productHandler.GetProductByID) // Get single product

	// Protected routes (vendor only)
	products.POST("", productHandler.CreateProduct, authMiddleware, loadUserMiddleware)       // Create product
	products.PUT("/:id", productHandler.UpdateProduct, authMiddleware, loadUserMiddleware)    // Update product
	products.DELETE("/:id", productHandler.DeleteProduct, authMiddleware, loadUserMiddleware) // Delete product

	// Vendor-specific routes
	vendor := g.Group("/vendor", authMiddleware, loadUserMiddleware)
	vendor.GET("/products", productHandler.GetVendorProducts) // Get my products
}

func setupShopRoutes(g *echo.Group, shopHandler *handler.ShopHandler, authMiddleware, loadUserMiddleware echo.MiddlewareFunc) {
	shops := g.Group("/shops")

	// Public routes
	shops.GET("", shopHandler.ListShops)                // List all shops
	shops.GET("/:id", shopHandler.GetShopByID)          // Get shop by ID
	shops.GET("/slug/:slug", shopHandler.GetShopBySlug) // Get shop by slug

	// Vendor routes (protected, vendor role required)
	vendorGroup := g.Group("", authMiddleware, loadUserMiddleware, middleware.RequireVendor())
	vendorGroup.POST("/shops", shopHandler.CreateShop)                   // Create shop
	vendorGroup.GET("/shops/my", shopHandler.GetMyShop)                  // Get my shop
	vendorGroup.GET("/shops/my/stats", shopHandler.GetMyShopStats)       // Get my shop stats
	vendorGroup.PUT("/shops/:id", shopHandler.UpdateShop)                // Update shop
	vendorGroup.PATCH("/shops/:id/status", shopHandler.ToggleShopStatus) // Toggle shop status

	// Admin routes (protected, admin role required)
	adminGroup := g.Group("", authMiddleware, loadUserMiddleware, middleware.RequireAdmin())
	adminGroup.DELETE("/shops/:id", shopHandler.DeleteShop)       // Delete shop
	adminGroup.PATCH("/shops/:id/verify", shopHandler.VerifyShop) // Verify shop
}

func setupCartRoutes(g *echo.Group, cartHandler *handler.CartHandler, authMiddleware, loadUserMiddleware echo.MiddlewareFunc) {
	cart := g.Group("/cart", authMiddleware, loadUserMiddleware)

	// Cart routes (all protected, customer role)
	cart.GET("", cartHandler.GetCart)                     // Get user's cart
	cart.GET("/count", cartHandler.GetCartItemCount)      // Get cart item count
	cart.POST("/items", cartHandler.AddToCart)            // Add item to cart
	cart.PUT("/items/:id", cartHandler.UpdateCartItem)    // Update cart item quantity
	cart.DELETE("/items/:id", cartHandler.RemoveCartItem) // Remove item from cart
	cart.DELETE("", cartHandler.ClearCart)                // Clear entire cart
}

func setupOrderRoutes(g *echo.Group, orderHandler *handler.OrderHandler, authMiddleware, loadUserMiddleware echo.MiddlewareFunc) {
	orders := g.Group("/orders", authMiddleware, loadUserMiddleware)

	// Customer order routes
	orders.POST("", orderHandler.CreateOrder)            // Create order (checkout)
	orders.GET("", orderHandler.GetOrders)               // Get user's orders
	orders.GET("/:id", orderHandler.GetOrderByID)        // Get order details
	orders.POST("/:id/cancel", orderHandler.CancelOrder) // Cancel order

	// Vendor order routes
	vendor := g.Group("/vendor", authMiddleware, loadUserMiddleware, middleware.RequireVendor())
	vendor.GET("/orders", orderHandler.GetVendorOrders)                // Get shop orders
	vendor.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus) // Update order status
}
