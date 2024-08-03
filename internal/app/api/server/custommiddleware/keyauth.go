package custommiddleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// KeyAuthMiddleware returns a KeyAuth middleware to check for access_token cookie.
func KeyAuthMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "cookie:access_token",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key != "", nil
		},
	})
}
