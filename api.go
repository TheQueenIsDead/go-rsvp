package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// InitAPI registers the application routes with the appropriate handlers,
// passing in a datastore wrapper for the endpoints to utilise.
func InitAPI(e *echo.Echo) {
	e.GET("/clicked", getClickedHandler)

	e.GET("/events", getEventsHandler)

}

// getClickedHandler returns hello world example text in order to fulfill
// the example clicked endpoint from the HTMX tutorial
func getClickedHandler(c echo.Context) error {

	log.Debug(c)
	log.Debug(c.Request().URL.Query())

	c.Response().WriteHeader(http.StatusOK)
	_, _ = c.Response().Write([]byte("Hello world!"))

	return nil
}

// getEventsHandler returns all events in the database
func getEventsHandler(c echo.Context) error {

	rows, err := RsvpDatabase.DB.Query("SELECT * FROM events")
	if err != nil {
		log.WithError(err).Error("could not retrieve events from database")
		c.Response().WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()

	var data []Events
	for rows.Next() {
		var e Events
		err = rows.Scan(&e.Id, &e.Time, &e.Description)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			c.Response().WriteHeader(http.StatusInternalServerError)
		}
		data = append(data, e)
	}

	var result string
	for _, e := range data {
		result += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", e.Time, e.Description)
	}

	c.Response().WriteHeader(http.StatusOK)
	_, _ = c.Response().Write([]byte(result))

	return nil
}
