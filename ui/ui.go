package ui

import (
	"fmt"
	"github.com/cbroglie/mustache"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/consts"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/models"
	"go-rsvp/templates"
	"google.golang.org/api/idtoken"
	"gorm.io/datatypes"
	"net/http"
	"strconv"
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
	//
	//// Partial Paths
	//app.Server.GET("/partial/events", eventsPartial)
	//app.Server.GET("/partial/events/:id", eventsByIdPartial)
	//app.Server.GET("/partial/events/new", eventsCreationPartial)
	//
	//// Components
	//app.Server.GET("/loginNavItem", loginNavItem)
	//
	//// Testing
	//app.Server.GET("/test", test)

}

func test(c echo.Context) error {

	//template := templates.Hello("hello")
	//return template.Render(c.Request().Context(), c.Response().Writer)
	//
	log.Debug(c.Request().Header)

	e := templates.Events([]models.Event{
		{
			Date:             models.EventDate{Date: datatypes.Date{}},
			Time:             models.EventTime{},
			Name:             "",
			Description:      "",
			MinimumAttendees: 0,
			Emoji:            "",
		},
	})

	return templates.Index(e).Render(c.Request().Context(), c.Response().Writer)

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

	el, err := database.GetEvents()
	if err != nil {
		log.WithError(err).Error("could not get events from database")
	}

	component := templates.Events(el)

	// Render the full page if the request was initiated by HTMX
	headers, ok := c.Request().Header["Hx-Request"]
	if ok && len(headers) == 1 && headers[0] == "true" {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Else full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)
}

// /events/:id
func eventsById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.WithError(err).WithField("event_id", c.Param("id")).Error("failed to parse event id")
	}

	event, err := database.GetEventById(id)
	if err != nil {
		log.WithError(err).Error("failed to parse event id")
	}

	attendees, err := database.GetAttendeesForEvent(event)
	if err != nil {
		log.WithError(err).Error("failed to retrieve attendees for event")
	}

	component := templates.Event(event, attendees)

	// Render the component if the request was initiated by HTMX
	headers, ok := c.Request().Header["Hx-Request"]
	if ok && len(headers) == 1 && headers[0] == "true" {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Else full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)

}

// /events/new
func eventsCreation(c echo.Context) error {

	component := templates.NewEvent()

	// Render the full page if the request was initiated by HTMX
	headers, ok := c.Request().Header["Hx-Request"]
	if ok && len(headers) == 1 && headers[0] == "true" {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)
}
