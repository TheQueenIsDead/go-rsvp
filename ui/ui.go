package ui

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/consts"
	"go-rsvp/container"
	"go-rsvp/models"
	"google.golang.org/api/idtoken"
	"net/http"
	"time"
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

	// Auth
	app.Server.GET("/login", login)

	// Paths
	app.Server.GET("/events", events)
	app.Server.GET("/events/:id", eventsById)
	app.Server.GET("/events/new", eventsCreation)

	// Partial Paths
	app.Server.GET("/partial/events", eventsPartial)
	app.Server.GET("/partial/events/:id", eventsByIdPartial)
	app.Server.GET("/partial/events/new", eventsCreationPartial)

	// Components
	app.Server.GET("/loginNavItem", loginNavItem)

}

func loginNavItem(c echo.Context) error {

	ctx := c.Request().Context()

	html := `<li style="float:right"><a class="active" href="/login">Logged out?! Mystery Man!!</a></li>`
	if cookie, _ := c.Request().Cookie("google"); cookie != nil {
		validate, err := idtoken.Validate(ctx, cookie.Value, consts.GoogleClientId)
		if err != nil {
			return err
		}

		name := validate.Claims["given_name"]
		imageUri := validate.Claims["picture"]
		log.Debug(imageUri)
		html = fmt.Sprintf(`<li style="float:right">
					<a class="active" href="/login">
						<img src="%s" referrerpolicy="no-referrer" class="rounded-circle" style="width: 25px" /> Log out as %s!
					</a>
				</li>`, imageUri, name)
	}
	return c.HTML(200, html)
}

func notFound(c echo.Context) error {
	return c.Render(200, "templates/template.404.html", nil)
}

// /login
// https://developers.google.com/identity/gsi/web/guides/get-google-api-clientid
func login(c echo.Context) error {
	output, err := mustache.RenderFile("templates/layout.login.html")
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}
	return c.HTML(200, output)
}

// /events
func events(c echo.Context) error {
	return c.Render(200, "templates/template.events.html", nil)
}
func eventsPartial(c echo.Context) error {
	output, err := mustache.RenderFile("templates/template.events.html")
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}
	return c.HTML(200, output)
}

// /events/:id
func eventsById(c echo.Context) error {
	return c.Render(200, "templates/template.event.html", map[string]interface{}{
		"eventId": c.Param("id"),
	})
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

// /events/new
func eventsCreation(c echo.Context) error {
	return c.Render(200, "templates/template.event.create.html", map[string]interface{}{
		"today": time.Now().Format(models.EventDateFormat),
	})
}
func eventsCreationPartial(c echo.Context) error {
	output, err := mustache.RenderFile("templates/template.event.create.html", map[string]interface{}{
		"today": time.Now().Format(models.EventDateFormat),
	})
	if err != nil {
		log.WithError(err).Error("could not render")
		return err
	}
	return c.HTML(200, output)
}
