package models

import "gorm.io/gorm"

type Attendee struct {
	gorm.Model
	Name    string
	EventId int
}
