package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/labstack/echo/v4"

	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/imbivek08/hamropasal/internal/model"
	"github.com/imbivek08/hamropasal/internal/service"
)

type contextKey string

const (
	ClerkUserIDKey contextKey = "clerk_user_id"
	ClerkEmailKey  contextKey = "clerk_email"
)

func ClerkAuthMiddleware(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
			}

			sessionToken := parts[1]

			claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{
				Token: sessionToken,
			})
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}

			userID := claims.Subject
			if userID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
			}

			ctx := context.WithValue(c.Request().Context(), ClerkUserIDKey, userID)
			c.SetRequest(c.Request().WithContext(ctx))

			c.Set("clerk_user_id", userID)

			return next(c)
		}
	}
}

func GetClerkUserID(c echo.Context) string {
	userID, _ := c.Get("clerk_user_id").(string)
	return userID
}

// LoadUserMiddleware loads the user from database and stores in context
func LoadUserMiddleware(userService *service.UserService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			clerkUserID := GetClerkUserID(c)
			if clerkUserID == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}

			user, err := userService.GetUserByClerkID(c.Request().Context(), clerkUserID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			c.Set("user", user)
			c.Set("user_role", user.Role)

			return next(c)
		}
	}
}

// RequireVendor middleware ensures the user is a vendor
func RequireVendor() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*model.User)
			if !ok || user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			if user.Role != "vendor" {
				return echo.NewHTTPError(http.StatusForbidden, "vendor access required")
			}

			return next(c)
		}
	}
}

// RequireAdmin middleware ensures the user is an admin
func RequireAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, ok := c.Get("user").(*model.User)
			if !ok || user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}

			if user.Role != "admin" {
				return echo.NewHTTPError(http.StatusForbidden, "admin access required")
			}

			return next(c)
		}
	}
}

func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok {
				return echo.NewHTTPError(http.StatusForbidden, "user role not found")
			}

			for _, role := range allowedRoles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(http.StatusForbidden, "insufficient permissions")
		}
	}
}

func InitClerkClient(cfg *config.Config) {
	clerk.SetKey(cfg.ClerkSecretKey)
}
