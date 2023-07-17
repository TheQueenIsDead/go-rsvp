package main

import (
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// InitAPI registers the application routes with the appropriate handlers,
// passing in a datastore wrapper for the endpoints to utilise.
func InitAPI(e *echo.Echo) {
	e.GET("/demo", getDemo)
	e.GET("/clicked", getClickedHandler)

	e.GET("/events", getEventsHandler)
	e.GET("/events/:id", getEventById)
	e.POST("/events/:id/attend", createEventAttendance)

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

	rows, err := RsvpDatabase.DB.Query("SELECT * FROM events")
	if err != nil {
		log.WithError(err).Error("could not retrieve events from database")
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()

	var data []Event
	for rows.Next() {
		var e Event
		err = rows.Scan(&e.Id, &e.Time, &e.Description)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		data = append(data, e)
	}

	//var result string
	//for _, e := range data {
	//	result += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", e.Time, e.Description)
	//}
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
		//result += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", e.Time, e.Description)
		result = append(result, evt)
	}

	return c.JSON(http.StatusOK, result)

}

// getEventsHandler returns an event by its ID
func getEventById(c echo.Context) error {

	id := c.Param("id")

	row := RsvpDatabase.DB.QueryRow("SELECT * FROM events WHERE id = ?", id)

	var e Event
	err := row.Scan(&e.Id, &e.Time, &e.Description)
	if err != nil {

		if err == sql.ErrNoRows {
			return c.Render(http.StatusNotFound, "templates/404.html", nil)
		}
		log.WithError(err).Error("could not unmarshal events from database")
		c.Response().WriteHeader(http.StatusInternalServerError)
	}

	log.Warn("Hello am in event by id")

	return c.Render(200, "templates/event.html", e)

}

// createEventAttendance allows us to register people for an event
// Ex. curl -v -X POST http://localhost:3000/events/40/attend -d 'name=billynomates2'
func createEventAttendance(c echo.Context) error {

	name := c.FormValue("name")
	id := c.Param("id")

	create := `insert into attendees (name, event_id) values (?, ?);`

	_, err := RsvpDatabase.DB.Exec(create, name, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "could not create attendee")
	}

	return c.String(200, fmt.Sprintf("All good for %s %s", name, id))

}

func getDemo(c echo.Context) error {

	return c.Render(200, "", nil)

	//content, err := mustache.RenderFileInLayout("templates/template.example.html", "templates/layout.index.html", nil)
	//if err != nil {
	//	log.Error(err)
	//}
	//c.Response().Header().Add("content-type", "text/html")
	//return c.String(200, content)
}
