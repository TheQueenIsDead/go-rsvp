package main

import (
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
)

var (
	RsvpDatabase *Database
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	output, err := mustache.RenderFileInLayout("templates/template.example.html", "templates/layout.index.html", nil)
	//output, err := mustache.RenderFile(name, data)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(output))
	if err != nil {
		return err
	}

	return nil
}

func main() {

	// Setup Logging
	log.StandardLogger().SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	// Setup Webserver
	e := echo.New()

	// Setup file server
	log.Info("initialising file server")
	e.Static("/", "static/")
	log.Info("fileserver initialised")
	//log.Info("fileserver initialised")

	// NOTE: https://echo.labstack.com/docs/templates
	e.Renderer = &Template{}

	//e.RouteNotFound("/*", func(c echo.Context) error {
	//	return c.Render(http.StatusNotFound, "templates/404.html", nil)
	//})

	//e.HTTPErrorHandler = customHTTPErrorHandler

	// Setup DB
	RsvpDatabase = NewDatabase()
	RsvpDatabase.Init()

	// Setup API
	log.Info("initializing api")
	InitAPI(e)
	log.Info("api initialised")

	// Serve
	e.Logger.Fatal(e.Start(":6969"))
}

//func customHTTPErrorHandler(err error, c echo.Context) {
//	httpError, ok := err.(*echo.HTTPError)
//	if ok {
//		errorCode := httpError.Code
//		switch errorCode {
//		case http.StatusNotFound:
//			// TODO: Smooth this over a bit
//			err := c.Redirect(http.StatusPermanentRedirect, "/404.html")
//			if err != nil {
//				return
//			}
//		default:
//			// TODO handle any other case
//			log.Debug("misc error thrown, yikes")
//		}
//	}
//
//}
