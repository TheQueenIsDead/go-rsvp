package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// InitAPI registers the application routes with the appropriate handlers,
// passing in a datastore wrapper for the endpoints to utilise.
func InitAPI() {
	http.HandleFunc("/clicked", getClickedHandler)
	http.HandleFunc("/events", getEventsHandler)
}

// getClickedHandler returns hello world example text in order to fulfill
// the example clicked endpoint from the HTMX tutorial
func getClickedHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte("Hello world!"))
}

// getEventsHandler returns all events in the database
func getEventsHandler(res http.ResponseWriter, req *http.Request) {

	rows, err := RsvpDatabase.DB.Query("SELECT * FROM events")
	if err != nil {
		log.WithError(err).Error("could not retrieve events from database")
		res.WriteHeader(http.StatusInternalServerError)
	}
	defer rows.Close()

	var data []Events
	for rows.Next() {
		var e Events
		err = rows.Scan(&e.Id, &e.Time, &e.Description)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			res.WriteHeader(http.StatusInternalServerError)
		}
		data = append(data, e)
	}

	var result string
	for _, e := range data {
		result += fmt.Sprintf("<tr><td>%s</td><td>%s</td></tr>", e.Time, e.Description)
	}

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write([]byte(result))
}
