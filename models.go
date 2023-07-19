package main

import (
	"time"
)

type Event struct {
	Id          int
	Time        time.Time
	Description string
}

type Attendee struct {
	Name    string
	EventId int
}
