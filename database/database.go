package database

import (
	log "github.com/sirupsen/logrus"
	"go-rsvp/container"
	"go-rsvp/models"
	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

const databaseFileName = "rsvp.db"

var (
	app container.Application
)

func NewDatabase() *gorm.DB {

	db, err := gorm.Open(sqlite.Open(databaseFileName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func Init(db *gorm.DB) {

	app.Database = db

	log.Info("initialising database")

	// Migrate the schema
	schemas := []interface{}{
		&models.Event{},
		&models.Attendee{},
	}
	for _, s := range schemas {
		err := db.AutoMigrate(s)
		if err != nil {
			log.WithError(err).WithField("schema", s).Panic("could not automigrate schema")
		}
	}

	// Create
	events := []models.Event{
		{Date: models.EventDate{Date: datatypes.Date(time.Now())}, Name: "Beers", Description: "Dizzy with the fizzy!", MinimumAttendees: 0},
		{Date: models.EventDate{Date: datatypes.Date(time.Now())}, Name: "Pool", Description: "Time to shark!", MinimumAttendees: 4},
		{Date: models.EventDate{Date: datatypes.Date(time.Now())}, Name: "Quiz", Description: "At the Dux!", MinimumAttendees: 6},
		{Date: models.EventDate{Date: datatypes.Date(time.Now())}, Name: "Puppy Walk", Description: "Heading to the New Brighton beach, usual place :^)", MinimumAttendees: 0},
	}
	for _, e := range events {
		res := db.Create(&e)
		err := res.Error
		if err != nil {
			log.WithError(err).Panic("could not insert events on startup")
		}
	}

	log.Info("database initialised")
}
