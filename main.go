package main

//go:generate: templ generate templates/

import (
	"errors"
	"go-rsvp/api"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/middleware"
	"go-rsvp/ui"
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
	e.Static("/", "static/")

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

	// Initialise Database
	db := database.NewDatabase()
	database.Init(db)

	// Initialise Routes
	app := container.Application{
		Server:   e,
		Database: db,
	}
	api.Init(app)
	ui.Init(app)

	// Serve
	e.Logger.Fatal(app.Server.Start(":3000"))
}
