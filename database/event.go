package database

import (
	log "github.com/sirupsen/logrus"
	"go-rsvp/models"
)

func getEvents() (data []models.Event, err error) {

	rows, err := app.Database.Query("SELECT * FROM events")
	defer rows.Close()
	if err != nil {
		log.WithError(err).Error("could not retrieve events from database")
		return
	}

	for rows.Next() {
		var e models.Event
		err = rows.Scan(&e.Id, &e.Time, &e.Description)
		if err != nil {
			log.WithError(err).Error("could not unmarshal events from database")
			return
		}
		data = append(data, e)
	}
	return
}
