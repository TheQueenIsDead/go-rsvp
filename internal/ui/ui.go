package ui

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/internal"
	"go-rsvp/internal/database"
	"go-rsvp/web/templates"
	"google.golang.org/api/idtoken"
	"net/http"
	"strconv"
)

var (
	db database.Database
)

func RegisterUIRoutes(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error { return c.Redirect(http.StatusPermanentRedirect, "/Events") })
	e.GET("/404", NotFound)
	e.GET("/login", Login)
	e.GET("/events", Events)
	e.GET("/events/:id", EventsById)
	e.GET("/events/new", EventsCreation)
}

func RegisterDatabase(d database.Database) {
	db = d
}

func loginNavItem(c echo.Context) error {

	ctx := c.Request().Context()

	html := `<li style="float:right"><a class="active" href="/Login">Logged out?! Mystery Man!!</a></li>`
	if cookie, _ := c.Request().Cookie("google"); cookie != nil {
		validate, err := idtoken.Validate(ctx, cookie.Value, internal.GoogleClientId)
		if err != nil {
			return err
		}

		name := validate.Claims["given_name"]
		imageUri := validate.Claims["picture"]
		log.Debug(imageUri)
		html = fmt.Sprintf(`<li style="float:right">
					<a class="active" href="/Login">
						<img src="%s" referrerpolicy="no-referrer" class="rounded-circle" style="width: 25px" /> Log out as %s!
					</a>
				</li>`, imageUri, name)
	}
	return c.HTML(200, html)
}

func NotFound(c echo.Context) error {
	component := templates.NotFound()
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)
}

// /Login
// https://developers.google.com/identity/gsi/web/guides/get-google-api-clientid
func Login(c echo.Context) error {

	component := templates.Login()

	// Else full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)

}

// /Events
func Events(c echo.Context) error {

	el, err := db.GetEvents()
	if err != nil {
		log.WithError(err).Error("could not get Events from database")
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

// /Events/:id
func EventsById(c echo.Context) error {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.WithError(err).WithField("event_id", c.Param("id")).Error("failed to parse event id")
	}

	event, err := db.GetEventById(id)
	if err != nil {
		log.WithError(err).Error("failed to parse event id")
	}

	attendees, err := db.GetAttendeesForEvent(event)
	if err != nil {
		log.WithError(err).Error("failed to retrieve attendees for event")
	}

	var attending bool
	for _, attendee := range attendees {
		userEmail, _ := c.Get("userEmail").(string)
		if userEmail == attendee.Email {
			attending = true
			break
		}
	}
	component := templates.Event(event, attendees, attending)

	// Render the component if the request was initiated by HTMX
	headers, ok := c.Request().Header["Hx-Request"]
	if ok && len(headers) == 1 && headers[0] == "true" {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Else full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)

}

// /Events/new
func EventsCreation(c echo.Context) error {

	component := templates.NewEvent()

	// Render the full page if the request was initiated by HTMX
	headers, ok := c.Request().Header["Hx-Request"]
	if ok && len(headers) == 1 && headers[0] == "true" {
		return component.Render(c.Request().Context(), c.Response().Writer)
	}

	// Full render
	return templates.Index(component).Render(c.Request().Context(), c.Response().Writer)
}
