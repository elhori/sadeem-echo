package infra

import "github.com/labstack/echo/v4"

func AuthorizationMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role").(string)
			for _, role := range roles {
				if userRole == role {
					return next(c)
				}
			}
			return echo.ErrForbidden
		}
	}
}

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.ErrUnauthorized
		}

		claims, err := ParseToken(token)
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Safely perform type assertion for user ID
		userID, ok := claims["id"].(float64)
		if !ok {
			return echo.ErrUnauthorized
		}

		// Pass userID to context for later use
		c.Set("user", int(userID))

		return next(c)
	}
}
