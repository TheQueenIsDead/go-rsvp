package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/database"
	"go-rsvp/models"
	"net/http"
	"strconv"
)

// GetClickedHandler returns hello world example text in order to fulfill
// the example clicked endpoint from the HTMX tutorial
func GetClickedHandler(c echo.Context) error {

	result := []map[string]interface{}{
		{
			"time":        "exampleTime",
			"description": "exampleDescription",
			"_links": map[string]string{
				"self": "/Events/demoId",
			},
		},
		{
			"time":        "exampleTime2",
			"description": "exampleDescription2",
			"_links": map[string]string{
				"self": "/Events/demoId2",
			},
		},
	}

	return c.JSON(http.StatusOK, result)

}

// getEventsHandler returns an event by its ID
func GetEventById(c echo.Context) error {

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

// CreateEventAttendance allows us to register people for an event
// Ex. curl -v -X POST http://localhost:3000/events/40/attend -d 'name=billynomates2'
func CreateEventAttendance(c echo.Context) error {

	userData := c.Get("userdata").(map[string]interface{})
	name, ok := userData["name"].(string)
	if !ok {
		log.Errorf("could not decode name from %v", userData)
	}

	email, ok := userData["email"].(string)
	if !ok {
		log.Errorf("could not decode email from %v", userData)
	}

	id, _ := strconv.Atoi(c.Param("id"))

	attendee := models.Attendee{
		Name:    name,
		Email:   email,
		EventId: id,
	}

	res := app.Database.Create(&attendee)
	err := res.Error
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not create attendee")

	}

	return c.String(200, fmt.Sprintf("All good for %s %d", name, id))

}

func CreateEvent(c echo.Context) error {

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
