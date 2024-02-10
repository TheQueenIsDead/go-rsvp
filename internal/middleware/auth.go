package middleware

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/internal"
	"google.golang.org/api/idtoken"
	"net/http"
	"time"
)

func CheckCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Path() == "/login" {
			return next(c)
		}

		cookie, err := c.Cookie("google")
		if err != nil || cookie == nil {
			_ = c.Redirect(http.StatusTemporaryRedirect, "/login")
			return next(c)
		}

		// Validate Cookie
		ctx := c.Request().Context()
		validator, err := idtoken.NewValidator(ctx)
		if err != nil {
			log.WithError(err).Error("could not create new google token validator")
			return c.Redirect(302, "/login")
		}
		validate, err := validator.Validate(ctx, cookie.Value, internal.GoogleClientId)
		if err != nil || validate == nil {
			log.WithError(err).Error("could not validate google id token")
			// Invalidate cookie
			c.SetCookie(&http.Cookie{
				Name:    "google",
				Value:   "invalid",
				Expires: time.Time{},
			})
			return c.Redirect(302, "/login")
		}

		// Propagate user information in ctx
		c.Set("userdata", validate.Claims)

		return next(c)

	}
}
