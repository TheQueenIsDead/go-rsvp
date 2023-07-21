package database

import (
	log "github.com/sirupsen/logrus"
	"go-rsvp/models"
)

func GetEvents() ([]models.Event, error) {
	var events []models.Event
	result := app.Database.Find(&events)
	err := result.Error
	if err != nil {
		log.WithError(err).Error("database: could not retrieve events")
	}
	return events, err
}

func GetEventById(id int) (models.Event, error) {
	var event models.Event
	result := app.Database.Find(&event, id)
	err := result.Error
	if err != nil {
		log.WithError(err).WithField("id", id).Error("database: could not retrieve event")
	}
	return event, err
}
