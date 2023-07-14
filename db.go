package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"time"
)

const databaseFileName = "rsvp.db"

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {

	log.Info("creating new database")
	db := &sql.DB{}
	db, err := sql.Open("sqlite3", databaseFileName)
	if err != nil {
		log.WithError(err).Panic("could not open database")
	}
	return &Database{DB: db}
}

func (d *Database) Init() {
	log.Info("initialising database")

	// Create events table
	create := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER NOT NULL PRIMARY KEY,
		time DATETIME NOT NULL,
		description TEXT
	);`
	res, err := d.DB.Exec(create)
	if err != nil {
		log.WithError(err).Panicf("could not init database: %s", res)
	}

	// Add example
	res, err = d.DB.Exec("INSERT INTO events VALUES(NULL,?,?);", time.Now(), "Beers")
	res, err = d.DB.Exec("INSERT INTO events VALUES(NULL,?,?);", time.Now(), "Pool")
	res, err = d.DB.Exec("INSERT INTO events VALUES(NULL,?,?);", time.Now(), "Quiz")
	res, err = d.DB.Exec("INSERT INTO events VALUES(NULL,?,?);", time.Now(), "Example event with a slightly longer description, something something coffee at New Brighton beach.")
	if err != nil {
		log.WithError(err).Panicf("could not init database: %s", res)
	}

	log.Infof("database initialised")
}
