package models

import (
	"time"
)

type Event struct {
	Id          int
	Time        time.Time
	Description string
}
