package infra

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func JWTAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the token from the Authorization header
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}

		// Check if the token has the correct format
		parts := strings.Split(token, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return echo.ErrUnauthorized
		}

		// Validate JWT token
		claims, err := ParseToken(parts[1])
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Extract user ID from claims
		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			return echo.ErrUnauthorized
		}
		userID := int(userIDFloat)

		// Set user ID in context
		c.Set("user", userID)

		// Call the next handler
		if err := next(c); err != nil {
			// Handle any errors returned by the next handler
			c.Error(err)
			return err
		}

		return nil
	}
}
