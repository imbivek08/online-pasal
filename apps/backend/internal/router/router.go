package router

import (
	"net/http"

	"github.com/imbivek08/hamropasal/internal/database"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, db *database.Database) {
	// Health check endpoint
	e.GET("/health", healthCheck)

	// API v1 group
	v1 := e.Group("/api/v1")

	// Example: User routes
	// userHandler := handler.NewUserHandler(service.NewUserService(db))
	// v1.POST("/users", userHandler.Create)
	// v1.GET("/users/:id", userHandler.GetByID)
	// v1.GET("/users", userHandler.List)
	// v1.PUT("/users/:id", userHandler.Update)
	// v1.DELETE("/users/:id", userHandler.Delete)

	// Add more route groups here as you build features
	setupUserRoutes(v1, db)
	setupProductRoutes(v1, db)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"message": "server is running",
	})
}

// Placeholder route groups - implement these as you add features
func setupUserRoutes(g *echo.Group, db *database.Database) {
	// users := g.Group("/users")
	// Add user-related routes here
}

func setupProductRoutes(g *echo.Group, db *database.Database) {
	// products := g.Group("/products")
	// Add product-related routes here
}
