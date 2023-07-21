package models

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const EventTimeFormat = "2023-07-21T11:11"

type EventTimestamp time.Time

type Event struct {
	gorm.Model
	Date             EventDate `form:"date"`
	Name             string    `form:"name"`
	Description      string    `form:"description"`
	MinimumAttendees int8      `form:"minimumAttendees"`
}

// EventDate is a custom type that extends the sql ORM library datatype for database operation compatibility,
// while allowing us to implement unmarshalling functionality as part of the echo web framework struct binding process.
type EventDate struct {
	datatypes.Date
}

// UnmarshalParam decodes and assigns a value from a form or query param.
// It is invoked as part of the web framework entity binding process, and is triggered by the inclusion of the `form`
// tag on an EventDate object
func (ed *EventDate) UnmarshalParam(src string) error {
	t, err := time.Parse("2006-01-02T15:04", src)
	if err != nil {
		log.WithError(err).WithField("source", src).Error("could not unmarshal eventdate param")
		return err
	}

	*ed = EventDate{datatypes.Date(t)}
	return nil
}
