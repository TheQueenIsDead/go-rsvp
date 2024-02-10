package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/consts"
	"go-rsvp/database"
	"go-rsvp/models"
	"go-rsvp/templates"
	"google.golang.org/api/idtoken"
	"gorm.io/datatypes"
	"strconv"
)

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

	html := `<li style="float:right"><a class="active" href="/Login">Logged out?! Mystery Man!!</a></li>`
	if cookie, _ := c.Request().Cookie("google"); cookie != nil {
		validate, err := idtoken.Validate(ctx, cookie.Value, consts.GoogleClientId)
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

	el, err := database.GetEvents()
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
