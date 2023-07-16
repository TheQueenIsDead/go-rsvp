package main

import (
	"fmt"
	"time"
)

type Events struct {
	Id          int
	Time        time.Time
	Description string
}

func (e *Events) ToString() string {
	return fmt.Sprintf("%d, %s, %s", e.Id, e.Time, e.Description)
}

func (e Events) String() string {
	return fmt.Sprintf("%d, %s, %s", e.Id, e.Time, e.Description)
}
