package main

import (
	"errors"
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	echoLog "github.com/labstack/gommon/log"
	log "github.com/sirupsen/logrus"
	"go-rsvp/api"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/middleware"
	"go-rsvp/ui"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {

	output, err := mustache.RenderFileInLayout(name, "templates/layout.index.html", data)
	//output, err := mustache.RenderFile(name, data)
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}

	_, err = w.Write([]byte(output))
	if err != nil {
		log.WithError(err).Error("could not write render")
		return err
	}

	return nil
}

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

	// NOTE: https://echo.labstack.com/docs/templates
	e.Renderer = &Template{}

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
