package infra

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}

		// Split token header, e.g., "Bearer tokenString"
		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}

		// Validate JWT token
		claims, err := ParseToken(parts[1])
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Set user ID in context
		userID := uint(claims["id"].(float64))
		c.Set("user", userID)

		return next(c)
	}
}

func AuthorizeRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role from context
			userRole := c.Get("role").(string)

			// Check if user role is authorized
			for _, role := range roles {
				if userRole == role {
					return next(c)
				}
			}

			return echo.ErrForbidden
		}
	}
}
