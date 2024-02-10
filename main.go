package main

//go:generate: templ generate templates/

import (
	"errors"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	log "github.com/sirupsen/logrus"
)

var (
	app container.Application
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
	database.InitialiseDatabase(db)

	// Initialise Routes
	app = container.Application{
		Server:   e,
		Database: db,
	}

	// Register HTTP routes

	//// API
	api := app.Server.Group("/api")
	api.GET("/clicked", GetClickedHandler)
	//api.GET("/Events", getEventsHandler)
	api.POST("/events/new", CreateEvent)
	api.GET("/events/:id", GetEventById)
	api.POST("/events/:id/attend", CreateEventAttendance)

	//// UI
	app.Server.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/Events") })
	app.Server.GET("/404", NotFound)
	app.Server.GET("/login", Login)
	app.Server.GET("/events", Events)
	app.Server.GET("/events/:id", EventsById)
	app.Server.GET("/events/new", EventsCreation)

	// Serve
	e.Logger.Fatal(app.Server.Start(":3000"))
}
