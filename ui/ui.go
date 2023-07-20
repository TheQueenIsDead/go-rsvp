package ui

import (
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/container"
	"net/http"
)

var (
	app container.Application
)

// Init registers the application UI routes with the appropriate rendering handlers.
func Init(a container.Application) {

	app = a

	// Top Level Redirect
	app.Server.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/events") })
	app.Server.GET("/404", notFound)

	// Paths
	app.Server.GET("/events", events)
	app.Server.GET("/events/:id", eventsById)

	// Partial Paths
	app.Server.GET("/partial/events", eventsPartial)
	app.Server.GET("/partial/events/:id", eventsByIdPartial)

}

func notFound(c echo.Context) error {
	return c.Render(200, "templates/template.404.html", nil)
}

func events(c echo.Context) error {
	return c.Render(200, "templates/template.events.html", nil)
}

func eventsById(c echo.Context) error {
	return c.Render(200, "templates/template.event.html", map[string]interface{}{
		"eventId": c.Param("id"),
	})
}

func eventsPartial(c echo.Context) error {
	output, err := mustache.RenderFile("templates/template.events.html")
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}
	return c.HTML(200, output)
}

func eventsByIdPartial(c echo.Context) error {
	output, err := mustache.RenderFile("templates/template.event.html", map[string]interface{}{
		"eventId": c.Param("id"),
	})
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}
	return c.HTML(200, output)
}
