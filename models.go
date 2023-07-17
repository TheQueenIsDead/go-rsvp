package main

import (
	"fmt"
	"time"
)

type Event struct {
	Id          int
	Time        time.Time
	Description string
}

func (e Event) ToString() string {
	return fmt.Sprintf("%d, %s, %s", e.Id, e.Time, e.Description)
}

type Attendee struct {
	EventId int
	Name    string
}

func (a Attendee) ToString() string {
	return fmt.Sprintf("%d, %s", a.EventId, a.Name)
}
