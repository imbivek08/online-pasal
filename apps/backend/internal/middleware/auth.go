package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/imbivek08/hamropasal/internal/config"
	"github.com/labstack/echo/v4"
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

func GetClerkUserID(c echo.Context) (string, bool) {
	userID, ok := c.Get("clerk_user_id").(string)
	return userID, ok
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
