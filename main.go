package main

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	RsvpDatabase *Database
)

func main() {

	log.StandardLogger().SetReportCaller(true)
	log.SetLevel(log.DebugLevel)

	// Setup file server
	log.Info("initialising file server")
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Info("fileserver initialised")
	//log.Info("initialising file server")
	//fs := http.FileServer(HTMLDir{"static/"})
	//http.Handle("/", http.StripPrefix("/", fs))
	//log.Info("fileserver initialised")

	// Setup DB
	RsvpDatabase = NewDatabase()
	RsvpDatabase.Init()

	// Setup API
	log.Info("initializing api")
	InitAPI()
	log.Info("api initialised")

	// Serve
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.WithError(err).Panic("issue encountered serving http")
	}
}
