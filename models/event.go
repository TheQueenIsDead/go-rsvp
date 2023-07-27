package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const (
	EventTimeFormat = "15:04"
	EventDateFormat = "2006-01-02"
)

type Event struct {
	gorm.Model
	Date             EventDate `form:"date"`
	Time             EventTime `form:"time"`
	Name             string    `form:"name"`
	Description      string    `form:"description"`
	MinimumAttendees int8      `form:"minimumAttendees"`
}

// EventDate is a custom type that extends the sql ORM library datatype for database operation compatibility,
// while allowing us to implement post body unmarshalling functionality as part of the echo web framework.
type EventDate struct {
	datatypes.Date
}

func (ed *EventDate) UnmarshalParam(param string) error {

	t, err := time.Parse(EventDateFormat, param)
	if err != nil {
		return err
	}

	ed.Date = datatypes.Date(t)

	return nil
}

// EventTime is a custom type that extends the sql ORM library datatype for database operation compatibility,
// while allowing us to implement post body unmarshalling functionality as part of the echo web framework.
type EventTime struct {
	datatypes.Time
}

func (et *EventTime) UnmarshalParam(param string) error {

	t, err := time.Parse(EventTimeFormat, param)
	if err != nil {
		return err
	}

	et.Time = datatypes.NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())

	return nil
}
