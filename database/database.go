package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"go-rsvp/container"
	"time"
)

const databaseFileName = "rsvp.db"

var (
	app container.Application
)

func NewDatabase() *sql.DB {

	log.Info("creating new database")
	db := &sql.DB{}
	db, err := sql.Open("sqlite3", databaseFileName)
	if err != nil {
		log.WithError(err).Panic("could not open database")
	}
	return db
}

func Init(a container.Application) {

	app = a

	log.Info("initialising database")

	// Create events table
	create := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER NOT NULL PRIMARY KEY,
		time DATETIME NOT NULL,
		name TEXT,
		description TEXT,
	  	minimumAttendees INTEGER DEFAULT 0
	);`
	res, err := app.Database.Exec(create)
	if err != nil {
		log.WithError(err).Panicf("could not init database: %s", res)
	}

	// Create attendees table
	create = `
	create table if not exists attendees (
	name     integer,
	event_id integer
		constraint attendees_events_id_fk
			references events
	);`
	res, err = app.Database.Exec(create)
	if err != nil {
		log.WithError(err).Panicf("could not init database: %s", res)
	}

	// Add example
	res, err = app.Database.Exec("INSERT INTO events VALUES(NULL,?,?,?,NULL);", time.Now(), "Beers", "desc1")
	res, err = app.Database.Exec("INSERT INTO events VALUES(NULL,?,?,?,NULL);", time.Now(), "Pool", "desc2")
	res, err = app.Database.Exec("INSERT INTO events VALUES(NULL,?,?,?,NULL);", time.Now(), "Quiz", "desc3")
	res, err = app.Database.Exec("INSERT INTO events VALUES(NULL,?,?,?,NULL);", time.Now(), "Puppy walk", "Example event with a slightly longer description, something something coffee at New Brighton beach.")
	if err != nil {
		log.WithError(err).Panicf("could not init database: %s", res)
	}

	log.Infof("database initialised")
}
