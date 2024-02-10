package database

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Attendee struct {
	gorm.Model
	Name    string
	Email   string
	EventId int
}

func (d *Database) GetAttendeesForEvent(e Event) ([]Attendee, error) {
	var attendees []Attendee
	result := d.Find(&attendees, "event_id = ?", e.ID)
	err := result.Error
	if err != nil {
		log.WithError(err).WithField("event_id", e.ID).Error("database: could not retrieve attendees for event")
	}
	return attendees, err
}
