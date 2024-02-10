package main

//go:generate: templ generate templates/

import (
	"errors"
	"go-rsvp/internal"
	"go-rsvp/internal/api"
	"go-rsvp/internal/database"
	"go-rsvp/internal/middleware"
	"go-rsvp/internal/ui"
	"net/http"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	log "github.com/sirupsen/logrus"
)

func main() {

	// Setup Logging
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	// Setup Webserver
	e := echo.New()
	e.Logger.SetLevel(echoLog.DEBUG)

	e.Use(middleware.CheckCookies)

	// NOTE: https://echo.labstack.com/docs/static-files
	e.Static("/", internal.StaticWebDirectory)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var httpError *echo.HTTPError
		ok := errors.As(err, &httpError)
		if ok {
			errorCode := httpError.Code
			switch errorCode {
			case http.StatusNotFound:
				err := c.Redirect(http.StatusTemporaryRedirect, "/404")
				if err != nil {
					log.WithError(err).Error("error redirecting 404")
				}
			default:
				// TODO handle any other case
				log.Infof("unhandled http error: %v", err)
			}
		}
	}

	db := database.NewDatabase()
	api.RegisterApiRoutes(e)
	api.RegisterDatabase(db)
	ui.RegisterUIRoutes(e)
	ui.RegisterDatabase(db)

	// Serve
	e.Logger.Fatal(e.Start(":3000"))
}
