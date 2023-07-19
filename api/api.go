package api

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"go-rsvp/container"
	"go-rsvp/models"
	"net/http"
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
	api.GET("/events/:id", getEventById)
	api.POST("/events/:id/attend", createEventAttendance)

}

// getClickedHandler returns hello world example text in order to fulfill
// the example clicked endpoint from the HTMX tutorial
func getClickedHandler(c echo.Context) error {

	result := []map[string]interface{}{
		{
			"time":        "exampletime",
			"description": "exampleDescription",
			"_links": map[string]string{
				"self": "/events/demoId",
			},
		},
		{
			"time":        "exampletime2",
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

	rows, err := app.Database.Query("SELECT * FROM events")
	if err != nil {
		log.WithError(err).Error("could not retrieve events from database")
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()

	var data []models.Event
	for rows.Next() {
		var e models.Event
		err = rows.Scan(&e.Id, &e.Time, &e.Description)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		data = append(data, e)
	}

	var result []map[string]interface{}
	for _, e := range data {
		evt := map[string]interface{}{
			"time":        e.Time,
			"description": e.Description,
			"_links": map[string]interface{}{
				"self": fmt.Sprintf("/events/%d", e.Id),
				"_templates": map[string]interface{}{
					fmt.Sprintf("/events/%d", e.Id): map[string]interface{}{
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

	id := c.Param("id")

	row := app.Database.QueryRow("SELECT * FROM events WHERE id = ?", id)

	var e models.Event
	err := row.Scan(&e.Id, &e.Time, &e.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Render(http.StatusNotFound, "templates/404.html", nil)
		}
		log.WithError(err).Error("could not unmarshal events from database")
		c.Response().WriteHeader(http.StatusInternalServerError)
	}

	rows, err := app.Database.Query("SELECT * FROM attendees WHERE event_id = ?", id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not get attendees")
	}
	var data []models.Attendee
	for rows.Next() {
		var a models.Attendee
		err = rows.Scan(&a.Name, &a.EventId)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		data = append(data, a)
	}

	log.Warn("Hello am in event by id")

	response := map[string]interface{}{
		"event":     e,
		"attendees": data,
	}

	return c.JSON(200, response)

	// TODO: Remove all UI render functions from the UI handlers file
	//return c.Render(200, "templates/event.html", e)

}

// createEventAttendance allows us to register people for an event
// Ex. curl -v -X POST http://localhost:3000/events/40/attend -d 'name=billynomates2'
func createEventAttendance(c echo.Context) error {

	name := c.FormValue("name")
	id := c.Param("id")

	create := `insert into attendees (name, event_id) values (?, ?);`

	_, err := app.Database.Exec(create, name, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not create attendee")
	}

	return c.String(200, fmt.Sprintf("All good for %s %s", name, id))

}
