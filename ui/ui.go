package ui

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
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
