package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

const EventTimeFormat = "2023-07-21T11:11"

type EventTimestamp time.Time

type Event struct {
	gorm.Model
	Date             datatypes.Date `form:"date"`
	Name             string         `form:"name"`
	Description      string         `form:"description"`
	MinimumAttendees int8           `form:"minimumAttendees"`
}

func (t *EventTimestamp) UnmarshalParam(src string) error {
	//ts, err := time.Parse(time.RFC3339, src)
	ts, err := time.Parse("2006-01-02T15:04", src)
	*t = EventTimestamp(ts)
	return err
}
