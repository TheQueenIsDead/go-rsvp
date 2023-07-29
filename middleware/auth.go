package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Process is the middleware function.
func CheckCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() == "/login" {
			return next(c)
		}

		// TODO: Actually propagate the JWT and verify
		// https://developers.google.com/identity/gsi/web/guides/verify-google-id-token
		_, err := c.Cookie("google")
		if err != nil {
			_ = c.Redirect(http.StatusTemporaryRedirect, "/login")
		}
		return next(c)

	}
}
