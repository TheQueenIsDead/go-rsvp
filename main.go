package main

import (
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

var (
	RsvpDatabase *Database
)

func main() {

	// Setup Logging
	log.StandardLogger().SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	// Setup Webserver
	e := echo.New()

	// Setup file server
	log.Info("initialising file server")
	e.Static("/", "static/")
	log.Info("fileserver initialised")
	//log.Info("fileserver initialised")

	// Setup DB
	RsvpDatabase = NewDatabase()
	RsvpDatabase.Init()

	// Setup API
	log.Info("initializing api")
	InitAPI(e)
	log.Info("api initialised")

	// Serve
	e.Logger.Fatal(e.Start(":3000"))
}
