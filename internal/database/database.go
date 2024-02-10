package database

import (
	log "github.com/sirupsen/logrus"
	"go-rsvp/internal"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	*gorm.DB
}

func NewDatabase() Database {

	db, err := gorm.Open(sqlite.Open(internal.DatabaseFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	database := Database{db}

	// Migrate the schema
	schemas := []interface{}{
		&Event{},
		&Attendee{},
	}
	for _, s := range schemas {
		err = database.AutoMigrate(s)
		if err != nil {
			log.WithError(err).WithField("schema", s).Panic("could not automigrate schema")
		}
	}

	log.Info("database initialised")

	return database
}

func createDemoEvents(db Database) error {
	// Create
	events := []Event{
		{Date: EventDate{datatypes.Date(time.Now())}, Name: "Beers", Description: "Dizzy with the fizzy!", MinimumAttendees: 0},
		{Date: EventDate{datatypes.Date(time.Now())}, Name: "Pool", Description: "Time to shark!", MinimumAttendees: 4},
		{Date: EventDate{datatypes.Date(time.Now())}, Name: "Quiz", Description: "At the Dux!", MinimumAttendees: 6},
		{Date: EventDate{datatypes.Date(time.Now())}, Name: "Puppy Walk", Description: "Heading to the New Brighton beach, usual place :^)", MinimumAttendees: 0},
	}
	var err error
	for _, e := range events {
		res := db.Create(&e)
		err = res.Error
		if err != nil {
			log.WithError(err).Error("could not insert demo events")
		}
	}
	return err
}
