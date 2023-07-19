package main

import (
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/api"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/ui"
	"html/template"
	"io"
	"net/http"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

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
	log.SetLevel(log.InfoLevel)

	// Setup Webserver
	e := echo.New()

	// NOTE: https://echo.labstack.com/docs/static-files
	e.Static("/", "static/")

	// NOTE: https://echo.labstack.com/docs/templates
	e.Renderer = &Template{}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		httpError, ok := err.(*echo.HTTPError)
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

	// Initialise Routes
	app := container.Application{
		Server:   e,
		Database: db,
	}
	database.Init(app)
	api.Init(app)
	ui.Init(app)

	// Serve
	e.Logger.Fatal(app.Server.Start(":3000"))
}
