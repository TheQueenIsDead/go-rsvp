package main

import (
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

// InitUI registers the application UI routes with the appropriate rendering handlers.
func InitUI(e *echo.Echo) {

	// Top Level Redirect
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/events") })
	e.GET("/404", notFound)

	// Paths
	e.GET("/events", events)
	e.GET("/events/:id", eventsById)

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
