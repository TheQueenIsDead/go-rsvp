package database

import (
	log "github.com/sirupsen/logrus"
	"go-rsvp/models"
)

func GetAttendeesForEvent(e models.Event) ([]models.Attendee, error) {
	var attendees []models.Attendee
	result := app.Database.Find(&attendees, "event_id = ?", e.ID)
	err := result.Error
	if err != nil {
		log.WithError(err).WithField("event_id", e.ID).Error("database: could not retrieve attendees for event")
	}
	return attendees, err
}
