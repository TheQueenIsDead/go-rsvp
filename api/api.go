package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/container"
	"go-rsvp/database"
	"go-rsvp/models"
	"net/http"
	"strconv"
)

var (
	app container.Application
)

// Init registers the application routes with the appropriate handlers,
// passing in a datastore wrapper for the endpoints to utilise.
func Init(a container.Application) {

	app = a

	api := app.Server.Group("/api")

	api.GET("/clicked", getClickedHandler)
	api.GET("/events", getEventsHandler)
	api.POST("/events/new", createEvent)
	api.GET("/events/:id", getEventById)
	api.POST("/events/:id/attend", createEventAttendance)

}

// getClickedHandler returns hello world example text in order to fulfill
// the example clicked endpoint from the HTMX tutorial
func getClickedHandler(c echo.Context) error {

	result := []map[string]interface{}{
		{
			"time":        "exampleTime",
			"description": "exampleDescription",
			"_links": map[string]string{
				"self": "/events/demoId",
			},
		},
		{
			"time":        "exampleTime2",
			"description": "exampleDescription2",
			"_links": map[string]string{
				"self": "/events/demoId2",
			},
		},
	}

	return c.JSON(http.StatusOK, result)

}

// getEventsHandler returns all events in the database
func getEventsHandler(c echo.Context) error {

	events, err := database.GetEvents()
	if err != nil {
		log.WithError(err).Error("api: could not get events from database")
	}

	var result []map[string]interface{}
	for _, e := range events {
		evt := map[string]interface{}{
			"name":        e.Name,
			"description": e.Description,
			"date":        e.Date.String(),
			"time":        e.Time.String(),
			"icon":        e.Emoji,
			"_links": map[string]interface{}{
				"self":    fmt.Sprintf("/events/%d", e.ID),
				"partial": fmt.Sprintf("/partial/events/%d", e.ID),
				// TODO: Re-evaluate how _templates is actually used (Not currently)
				"_templates": map[string]interface{}{
					fmt.Sprintf("/events/%d", e.ID): map[string]interface{}{
						"title":       "Attend",
						"method":      "POST",
						"contentType": "application/json",
						"properties": []map[string]interface{}{
							{
								"name":     "name",
								"required": true,
							},
						},
					},
				},
			},
		}
		result = append(result, evt)
	}

	return c.JSON(http.StatusOK, result)

}

// getEventsHandler returns an event by its ID
func getEventById(c echo.Context) error {

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

	// TODO: Tidy this up a wee bit.
	response := map[string]interface{}{
		"event":     event,
		"attendees": attendees,
	}

	return c.JSON(200, response)
}

// createEventAttendance allows us to register people for an event
// Ex. curl -v -X POST http://localhost:3000/events/40/attend -d 'name=billynomates2'
func createEventAttendance(c echo.Context) error {

	name := c.FormValue("name")
	id, _ := strconv.Atoi(c.Param("id"))

	attendee := models.Attendee{
		Name:    name,
		EventId: id,
	}

	res := app.Database.Create(&attendee)
	err := res.Error
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not create attendee")

	}

	return c.String(200, fmt.Sprintf("All good for %s %d", name, id))

}

func createEvent(c echo.Context) error {

	var event models.Event
	err := c.Bind(&event)

	if err != nil {
		log.WithError(err).Error("could not bind")
		return c.String(http.StatusBadRequest, "bad request")
	}

	res := app.Database.Create(&event)
	err = res.Error
	if err != nil {
		log.WithError(err).Error("could not create event")
		return c.String(http.StatusInternalServerError, "could not process event")
	}

	return c.String(http.StatusOK, "OK")
}
