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
