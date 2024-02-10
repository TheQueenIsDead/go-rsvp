package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/internal/database"
	"net/http"
	"strconv"
)

var (
	db database.Database
)

func RegisterApiRoutes(e *echo.Echo) {
	//// API
	api := e.Group("/api")
	api.POST("/events/new", CreateEvent)
	api.POST("/events/:id/attend", CreateEventAttendance)
}

func RegisterDatabase(d database.Database) {
	db = d
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

	attendee := database.Attendee{
		Name:    name,
		Email:   email,
		EventId: id,
	}

	res := db.Create(&attendee)
	err := res.Error
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not create attendee")

	}

	return c.String(200, fmt.Sprintf("All good for %s %d", name, id))

}

func CreateEvent(c echo.Context) error {

	var event database.Event
	err := c.Bind(&event)

	if err != nil {
		log.WithError(err).Error("could not bind")
		return c.String(http.StatusBadRequest, "bad request")
	}

	res := db.Create(&event)
	err = res.Error
	if err != nil {
		log.WithError(err).Error("could not create event")
		return c.String(http.StatusInternalServerError, "could not process event")
	}

	return c.String(http.StatusOK, "OK")
}
